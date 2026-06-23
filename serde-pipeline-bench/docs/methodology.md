# Methodology

We compare JSON, Protobuf, and Avro on the same record, measuring wire size and
encode/decode cost, then confirm it end-to-end in a resource-capped pipeline.

## The formats and what differs

| Format   | Encoding | Field identity on wire | Numbers          |
| -------- | -------- | ---------------------- | ---------------- |
| json     | text     | full field names       | ASCII            |
| protobuf | binary   | numeric field tags     | varint / fixed64 |
| avro     | binary   | none (positional)      | zigzag / fixed   |

JSON differs from the other two on **two** axes simultaneously: it's text (so it
pays character scanning + ASCII↔number conversion + escaping) and it carries full
field names. Protobuf and Avro are both binary and both drop the names. That's
why these three formats can quantify the *total* gap but not cleanly separate
"because binary" from "because no keys" — both change together. We call that out
rather than overclaim.

## Fairness rules

- **One data model.** Every codec serializes the same `event.Event` — 12 fields,
  a realistic mix of strings/floats/ints with normal (not artificially short)
  names. Only the encoding differs.
- **Deterministic data.** `event.Sample(seed)` is reproducible; the microbench
  rotates over 256 distinct records so we don't measure one cache-hot payload.
- **Correctness gate.** `TestRoundTrip` requires every codec to reproduce the
  record exactly before any timing is trusted.
- **Protobuf via `protowire`.** Hand-rolled encode/decode over the official
  low-level package is byte-identical to generated code and keeps the repo
  buildable without `protoc`.

## Microbenchmark

`go test ./bench -bench . -benchmem` reports ns/op + allocs/op for Marshal,
Unmarshal, and RoundTrip. Decode is weighted most heavily when interpreting,
because an ingest pipeline decodes far more than it encodes.

## Pipeline: holding Kafka lag constant

The end-to-end check confirms the microbench result isn't an artifact of a tight
in-cache loop.

**Constant-lag design.** Lag (unconsumed backlog) is a confound: if it varied per
run, a faster-looking format might just have had less to do. So we fix it:

1. Producer (uncapped — not under test) writes **exactly `N`** records to a fresh
   single-partition topic and exits. Starting state for every scenario: offset 0,
   lag `N`.
2. Consumer (CPU/memory **capped** via compose `deploy.resources.limits`) drains
   from the start until it has seen `N`, and reports wall-clock RPS. Timing starts
   on the first record, excluding connection setup.

Same `N`, same cap, same topic shape, same data — only the wire format changes.

**Knobs:** `N` (backlog/lag), `CONSUMER_CPUS`, `CONSUMER_MEM`.

## Known measurement artifacts (fix or disclose before publishing)

These are properties of our *code/libraries*, not the formats — do not quote them
as format characteristics without a caveat:

- **Protobuf encode allocates ~5×/248 B.** Our `Marshal` appends to a `nil` slice
  that reallocates as it grows. Generated protobuf sizes the buffer once
  (`proto.Size`) and allocates once. So our protobuf *encode* is penalized by
  buffer growth, not by the format. Fix: presize the buffer.
- **Avro decode shows 0 allocs / 36 B.** `hamba/avro` very likely points string
  fields *into the read buffer* (aliasing) instead of copying. Our Protobuf decode
  does `string(v)` per string field → 4 real copies → 4 allocs. That's not
  apples-to-apples: aliasing is faster but can be a correctness hazard if the
  source buffer is reused (as Kafka client buffers often are). Decide one policy
  (copy vs alias) for all codecs and state it.
- **Hand-rolled vs generated Protobuf.** Wire-identical, but a generated decoder
  has tuned field dispatch; ours is a `switch` in a loop. Likely close; note it.

## Threats to validity (call these out in the post)

- **Library quality ≠ format.** A slow JSON or Avro library would shift numbers.
- **Schema source.** Avro/Protobuf assume the schema is in hand. Real pipelines
  often fetch it from a schema registry — a cost this harness excludes. Avro
  *requires* the schema to read at all; Protobuf can partially self-describe.
- **Field shape matters.** Many long strings favor text less; many numeric fields
  favor binary more. Report the record shape alongside the numbers.
- **Single partition / single node.** Isolates decode cost; real throughput also
  depends on partitions, network, and the downstream sink.
