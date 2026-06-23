# serde-pipeline-bench

A controlled benchmark answering one question:

> **Are Protobuf pipelines much faster than JSON pipelines *because* the keys
> are dropped from the payload (defined and ordered in the schema instead)?**

Short answer from the data so far: **partly — but it's "binary *and* keyless,"
not keys alone.** This repo is built to *decompose* the win into its parts rather
than just declare a winner.

Linear: [SUD-5](https://linear.app/sudopower/issue/SUD-5).

## The hypothesis, and how we isolate it

"Protobuf has no keys on the wire" conflates two independent variables:

| Format    | Wire   | Keys on wire        | What it isolates                       |
| --------- | ------ | ------------------- | -------------------------------------- |
| json      | text   | yes (full names)    | baseline                               |
| jsonshort | text   | yes (1-char names)  | **key *length* alone** (still text)    |
| msgpack   | binary | yes (full names)    | **the control: binary but keyed**      |
| cbor      | binary | yes (full names)    | second binary-keyed control            |
| protobuf  | binary | no (field tags)     | keyless binary                         |
| avro      | binary | no (positional)     | keyless binary, no tags at all         |

The logic:

- **json → jsonshort** isolates how much is just *long field names* in text.
- **json → msgpack/cbor** isolates *going binary while keeping keys*.
- **msgpack/cbor → protobuf/avro** isolates *dropping the keys*.

If "no keys" were the whole story, msgpack/cbor (binary **with** keys) would sit
close to JSON. If they sit close to protobuf/avro, the win is *binary*, not keys.

## What's here

- `internal/event` — one canonical record shared by every format (so only the
  serialization varies).
- `internal/codec` — one file per format behind a common interface, plus a
  round-trip correctness test (`go test ./internal/codec`).
- `bench/` — microbenchmarks: marshal / unmarshal / round-trip, with allocs.
- `cmd/sizes` — wire-size report (size is the variable the hypothesis is about).
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
# all formats, 2M msgs each, consumer capped at 1 CPU / 512M
./run-pipeline.sh

# tweak the controls
N=1000000 CONSUMER_CPUS=0.5 CONSUMER_MEM=256M ./run-pipeline.sh json protobuf
```

Results land in `pipeline-results.txt` (grep `RESULT`).

## Findings

See [`docs/results.md`](docs/results.md) for the numbers and the running
interpretation. [`docs/methodology.md`](docs/methodology.md) documents every
control and the reasoning, so this repo doubles as the draft spine of the post.
