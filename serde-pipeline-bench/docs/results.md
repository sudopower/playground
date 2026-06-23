# Results

Machine: Apple Silicon (arm64), Go 1.25.3. Microbenchmark = `go test ./bench
-bench . -benchmem -benchtime 1s`, 256-record rotation. Pipeline numbers TBD.

> These are local microbench numbers — absolute values vary by machine; the
> **ratios** are the story.

## Wire size (avg bytes/msg, 1000 records)

| Format    | Bytes | vs JSON | Keys on wire       |
| --------- | ----- | ------- | ------------------ |
| avro      | 92.6  | 30%     | none (positional)  |
| protobuf  | 104.1 | 34%     | numeric tags       |
| cbor      | 216.7 | 70%     | full names         |
| jsonshort | 216.8 | 70%     | 1-char names       |
| msgpack   | 226.2 | 73%     | full names         |
| json      | 309.8 | 100%    | full names         |

## Decode — ns/op (the number that matters most for ingest)

| Format    | ns/op | allocs/op | Wire   | Keys |
| --------- | ----- | --------- | ------ | ---- |
| protobuf  | 126   | 4         | binary | no   |
| avro      | 141   | 0         | binary | no   |
| msgpack   | 663   | 5         | binary | yes  |
| cbor      | 700   | 4         | binary | yes  |
| jsonshort | 2124  | 9         | text   | tiny |
| json      | 2452  | 8         | text   | full |

## Encode — ns/op

| Format    | ns/op | allocs/op |
| --------- | ----- | --------- |
| avro      | 132   | 1         |
| protobuf  | 140   | 5         |
| cbor      | 261   | 1         |
| msgpack   | 396   | 4         |
| json      | 646   | 1         |
| jsonshort | 677   | 2         |

---

## Decomposition — answering the hypothesis

Walking the controls for **decode**:

| Step                                  | Comparison         | Speedup | Attributed to        |
| ------------------------------------- | ------------------ | ------- | -------------------- |
| shorten keys, stay text               | json → jsonshort   | ~1.15×  | key *length*         |
| go binary, **keep** keys              | json → msgpack     | ~3.7×   | text → binary        |
| **drop** keys, stay binary            | msgpack → protobuf | ~5.3×   | dropping keys        |
| **total**                             | json → protobuf    | ~19×    | everything combined  |

### Verdict on "protobuf is fast *because* it drops keys"

**Substantially supported for decode — but incomplete.** Two findings:

1. **Dropping keys is a real, large lever** (~5.3× on decode here). It's not just
   smaller bytes — removing names lets the decoder go *positional/tag-directed*
   straight into fixed struct fields instead of matching key strings. So the
   hypothesis is more right than "it's just binary" skeptics assume.

2. **But going binary is itself ~3.7×, with keys still present.** msgpack and
   cbor — binary *with* full field names — land ~3.5–4× faster than JSON, nowhere
   near it. So roughly: *binary* and *keyless* each contribute a big multiplier,
   and they stack. It's not keys *alone*.

3. **The cleanest surprise: key length ≠ speed.** `jsonshort` reclaims most of
   the *size* (70% of JSON, same as binary-keyed) yet is only ~15% faster to
   decode. So in text formats, long names cost **bytes** but the dominant **CPU**
   cost is text parsing (number/string scanning), not the key strings.

**One-line takeaway for the post:** Protobuf's decode win is "binary × keyless,"
roughly 4× from leaving text behind and another ~5× from dropping field names so
the decoder can go positional. The popular "no keys in the payload" explanation
captures the bigger of the two factors but misses that half the win is simply not
being text.

## Pipeline (end-to-end, capped consumer, constant lag)

_Pending — run `./run-pipeline.sh` and paste `RESULT` lines here._

| Format | N | cap | elapsed | RPS | ns/msg |
| ------ | - | --- | ------- | --- | ------ |
| _tbd_  |   |     |         |     |        |

Expectation: the microbench ranking should hold, but absolute gaps **compress**
because Kafka fetch, network, and framing overhead are shared fixed costs that
don't care about the codec. How much they compress is itself a finding worth
reporting.
