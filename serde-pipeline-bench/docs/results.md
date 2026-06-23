# Results

Machine: Apple M2 Pro, Go 1.25.3. Microbenchmark = `go test ./bench -bench .
-benchmem -benchtime 1s`, 256-record rotation. Pipeline numbers TBD.

> Local microbench numbers — absolute values vary by machine; the **ratios** are
> the story. See `methodology.md` for measurement artifacts before quoting these.

## Wire size (avg bytes/msg, 1000 records)

| Format   | Bytes | vs JSON | Field identity |
| -------- | ----- | ------- | -------------- |
| avro     | 92.6  | 30%     | none (positional) |
| protobuf | 104.1 | 34%     | numeric tags   |
| json     | 309.8 | 100%    | full names     |

Avro edges Protobuf because Protobuf still spends ~1 tag byte per field (12
fields ≈ the 11.5-byte gap); Avro spends zero. That extra byte is what buys
Protobuf schema-evolution safety — Avro's positional layout has none.

## Decode — ns/op (the number that matters most for ingest)

| Format   | ns/op | allocs/op | vs JSON |
| -------- | ----- | --------- | ------- |
| protobuf | 123.7 | 4         | ~20×    |
| avro     | 138.0 | 0*        | ~18×    |
| json     | 2416  | 8         | 1×      |

\* See artifacts: Avro's 0 allocs is almost certainly string aliasing into the
read buffer; Protobuf's 4 are real `string()` copies.

## Encode — ns/op

| Format   | ns/op | allocs/op |
| -------- | ----- | --------- |
| avro     | 129.1 | 1         |
| protobuf | 136.6 | 5†        |
| json     | 634.8 | 1         |

† Artifact: our Protobuf `Marshal` grows a `nil` slice (5 allocs); a presized
buffer would be ~1. Not a format property.

## Round-trip — ns/op

| Format   | ns/op | allocs/op |
| -------- | ----- | --------- |
| protobuf | 264.0 | 9         |
| avro     | 275.8 | 1         |
| json     | 3099  | 9         |

## Interpretation

- **The binary/schema formats are ~18–20× faster to decode than JSON**, and
  decode is the half an ingest pipeline does most. JSON's decode (2416 ns) is ~4×
  its own encode (635 ns); Protobuf is roughly symmetric (124 vs 137). So JSON's
  worst case is the one you hit most under read-heavy load.
- **Size: ~3× smaller** for both binary formats (30–34% of JSON).
- **Why (qualitative, since these three formats can't isolate it):** going binary
  removes text parsing (ASCII number conversion, escapes, structural scanning),
  and dropping field names lets a schema-aware decoder write straight into fixed
  struct fields instead of matching key strings. Both contribute; with only
  json/protobuf/avro we can't put a clean percentage on each.
- **avro vs protobuf** is close on speed; Avro is smaller and protobuf carries
  tags (the price of forward/backward compatibility).

## Pipeline (end-to-end, capped consumer, constant lag)

N = 2,000,000 per format · consumer capped at 1.0 CPU / 512M · single partition ·
constant starting lag of 2M (producer fills topic fully, then consumer drains).

| Format   | elapsed | RPS   | ns/msg | vs JSON |
| -------- | ------- | ----- | ------ | ------- |
| avro     | 0.88s   | 2.27M | 440    | 7.1×    |
| protobuf | 0.89s   | 2.24M | 446    | 7.0×    |
| json     | 6.22s   | 321k  | 3111   | 1×      |

### The microbench gap compresses end-to-end (~20× → ~7×)

Decode in isolation was ~20× (protobuf 124 ns vs json 2416 ns). In the pipeline
it's **~7×**. The difference is fixed per-record overhead — Kafka fetch, record
iteration, framing — that every format pays regardless of codec:

- protobuf pipeline 446 ns/msg ≈ 124 ns decode **+ ~320 ns fixed overhead**.
  Decode is the minority of its time, so the overhead dominates and dilutes the
  win.
- json pipeline 3111 ns/msg ≈ 2416 ns decode + similar overhead. Decode still
  dominates, so the overhead barely moves it.

**Takeaway for the post:** the serialization choice is worth ~7× of *consumer
throughput* here, not the ~20× a microbenchmark alone would suggest. Quote the
end-to-end number, not the microbench, when talking about pipeline impact — and
note it scales with how decode-bound the consumer is (a heavier sink or transform
downstream would compress the gap further).

Avro and protobuf are a dead heat; Avro's marginally smaller payload (less to
fetch) gives it a hair's edge. Both are nowhere near JSON.

_Reproduce: `N=2000000 ./run-pipeline.sh json protobuf avro` (Apple M2 Pro, Go 1.25.3)._
