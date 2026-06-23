# Methodology

The goal is not to crown a format but to **attribute** the JSON→Protobuf speedup
to specific causes. That requires controls, not just contenders.

## The variables

Serialization formats differ along (at least) two axes that usually get
conflated when people say "protobuf is faster because it has no keys":

1. **Encoding** — text (JSON) vs binary (everything else). Text means
   character-by-character scanning, escape handling, and ASCII↔number conversion
   (`strconv`). Binary skips all of it.
2. **Keys on the wire** — does each record carry its field names?
   - JSON / MessagePack / CBOR: yes, full field names.
   - Protobuf: no names, just numeric field tags.
   - Avro: nothing at all — purely positional, order fixed by the schema.

A third axis falls out of the decoder design:

3. **Allocations / reflection** — generic JSON decode builds maps and leans on
   reflection (heavy GC pressure); schema-aware binary decoders write into fixed
   fields. We capture this with `-benchmem` (allocs/op, B/op).

## The controls

| Comparison              | Holds fixed        | Isolates            |
| ----------------------- | ------------------ | ------------------- |
| json → jsonshort        | text, reflection   | key **length**      |
| json → msgpack/cbor     | keys present       | text → **binary**   |
| msgpack/cbor → protobuf | binary             | **dropping keys**   |
| protobuf → avro         | binary, keyless    | tags vs positional  |

`msgpack` and `cbor` are the load-bearing controls: binary **but keyed**. If
dropping keys were the dominant factor, they'd cluster near JSON. Two independent
keyed-binary libraries guard against a single library's quirk skewing the result.

## Fairness rules

- **One data model.** Every codec serializes the same `event.Event`; only the
  encoding differs. Same 12 fields, realistic (not artificially short) names.
- **Deterministic data.** `event.Sample(seed)` is reproducible; the microbench
  rotates over 256 distinct records so we don't measure an unrealistically
  cache-friendly single payload.
- **Correctness gate.** `TestRoundTrip` requires every codec to reproduce the
  record exactly before any number is trusted.
- **Protobuf via `protowire`.** Hand-rolled encode/decode over the official
  low-level package produces byte-identical output to generated code (which calls
  the same package), and keeps the repo buildable without `protoc`.

## Microbenchmark

`go test ./bench -bench . -benchmem`. Reports ns/op + allocs/op for Marshal,
Unmarshal, and RoundTrip. Decode is weighted most heavily in interpretation
because an ingest pipeline decodes far more than it encodes.

## Pipeline: holding Kafka lag constant

The end-to-end check exists to confirm the microbench result isn't an artifact
of running in a tight loop with everything in cache.

**Constant-lag design.** "Lag" (unconsumed backlog) is a confound: if it varied
per run, a faster-looking format might just have had less to do. So we fix it:

1. Producer (uncapped — not under test) writes **exactly `N`** records to a fresh
   single-partition topic and exits. Starting state for every scenario: offset 0,
   lag `N`.
2. Consumer (CPU/memory **capped** via compose `deploy.resources.limits`) starts
   from the beginning, decodes until it has seen `N`, and reports wall-clock RPS.
   Timing starts on the first record, excluding connection setup.

Same `N`, same cap, same topic shape, same data — only the wire format changes,
so RPS is directly attributable to the codec.

**Knobs:** `N` (backlog/lag), `CONSUMER_CPUS`, `CONSUMER_MEM`.

## Threats to validity (call these out in the post)

- **Library quality ≠ format.** `hamba/avro` and our `protowire` path are fast;
  a slow JSON or Avro library would shift numbers. We note which libs are used.
- **Schema source.** Avro/Protobuf assume the schema is already in hand. A real
  pipeline often fetches it from a schema registry — a cost this harness excludes.
- **Field shape matters.** Many long string fields favor text less; many numeric
  fields favor binary more. Our record is a deliberate mix; results will shift
  with payload shape, so report the shape alongside the numbers.
- **Single partition / single node.** Isolates decode cost; real throughput also
  depends on partitions, network, and the downstream sink.
