# serde-pipeline-bench

A controlled benchmark for [SUD-5](https://linear.app/sudopower/issue/SUD-5)
comparing **JSON, Protobuf, and Avro** pipelines:

> How much faster — and smaller — are the schema-based binary formats (Protobuf,
> Avro) than JSON, and *why*?

The original hypothesis was that Protobuf wins mainly because keys are dropped
from the payload (defined and ordered in the schema instead). With these three
formats we can measure the **size** of the win precisely; the "why" is
qualitative (see [Caveats](#caveats)).

## The three formats

| Format   | Wire   | Field identity on the wire | Schema needed to read? |
| -------- | ------ | -------------------------- | ---------------------- |
| json     | text   | full field names           | no                     |
| protobuf | binary | numeric field tags         | (helps, not required)  |
| avro     | binary | nothing (positional)       | **yes**                |

The progression is "how much per-field identification travels on the wire":
full names → numeric tags → nothing. That maps directly onto both size and speed.

## What's here

- `internal/event` — one canonical record shared by every format (so only the
  serialization varies).
- `internal/codec` — one file per format behind a common interface, plus a
  round-trip correctness test (`go test ./internal/codec`).
- `bench/` — microbenchmarks: marshal / unmarshal / round-trip, with allocs.
- `cmd/sizes` — wire-size report.
- `pipeline/` — a local Kafka + **resource-capped** consumer to confirm the
  microbench result survives end-to-end under a **constant Kafka lag**.

Protobuf is encoded via the official low-level `protowire` package, so there is
**no `protoc` codegen step** — `go build ./...` is all you need.

## Run it

### Microbenchmarks (no infra)

```bash
go test ./internal/codec -run TestRoundTrip -v   # correctness first
go run ./cmd/sizes                               # wire sizes
go test ./bench -run x -bench . -benchmem        # speed + allocs
```

### Pipeline (Docker)

Constant-lag model: the producer fills a topic with exactly `N` messages **before**
the consumer starts, so every scenario drains the *same fixed backlog* under the
*same CPU/memory cap*. Throughput is therefore directly comparable across formats.

```bash
./run-pipeline.sh                                   # json, protobuf, avro
N=1000000 CONSUMER_CPUS=0.5 ./run-pipeline.sh json protobuf
```

Results land in `pipeline-results.txt` (grep `RESULT`).

## Findings

See [`docs/results.md`](docs/results.md) for numbers and interpretation, and
[`docs/methodology.md`](docs/methodology.md) for the controls and the honest
list of what this setup can and can't prove.

## Caveats

With only json / protobuf / avro, both binary formats change **two** things at
once versus JSON — they're binary *and* they drop field names — so the raw
numbers show *that* they're far faster, not a clean split of *how much* is "binary"
vs "no keys". The reasoning for the split is qualitative here. See
`docs/methodology.md` for the known measurement artifacts (e.g. our Protobuf
encode buffer, Avro string aliasing) before quoting any number in the post.
