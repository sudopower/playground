package bench

import (
	"testing"

	"serdebench/internal/codec"
	"serdebench/internal/event"
)

// sampleSet is a small pool of distinct records so the benchmark doesn't measure
// the same bytes over and over (which would be unrealistically cache-friendly).
var sampleSet = event.Batch(256)

// BenchmarkMarshal measures encode cost per format. Run with -benchmem to get
// allocs/op and B/op — allocations are usually where JSON loses, so they matter
// as much as ns/op for the story.
func BenchmarkMarshal(b *testing.B) {
	for _, c := range codec.All() {
		c := c
		b.Run(c.Name(), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				e := &sampleSet[i%len(sampleSet)]
				if _, err := c.Marshal(e); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkUnmarshal measures decode cost — the half that dominates an ingest
// pipeline, since the ingestor decodes far more than it encodes.
func BenchmarkUnmarshal(b *testing.B) {
	for _, c := range codec.All() {
		c := c
		// Pre-encode the pool once so we only time decoding.
		encoded := make([][]byte, len(sampleSet))
		for i := range sampleSet {
			data, err := c.Marshal(&sampleSet[i])
			if err != nil {
				b.Fatal(err)
			}
			encoded[i] = data
		}
		b.Run(c.Name(), func(b *testing.B) {
			b.ReportAllocs()
			var out event.Event
			for i := 0; i < b.N; i++ {
				if err := c.Unmarshal(encoded[i%len(encoded)], &out); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkRoundTrip measures encode+decode together — closest single number to
// "what one message costs as it passes through a stage".
func BenchmarkRoundTrip(b *testing.B) {
	for _, c := range codec.All() {
		c := c
		b.Run(c.Name(), func(b *testing.B) {
			b.ReportAllocs()
			var out event.Event
			for i := 0; i < b.N; i++ {
				data, err := c.Marshal(&sampleSet[i%len(sampleSet)])
				if err != nil {
					b.Fatal(err)
				}
				if err := c.Unmarshal(data, &out); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
