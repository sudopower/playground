// Command consumer drains a pre-filled topic, decoding every record with the
// chosen format, and reports end-to-end throughput. It is meant to run under a
// fixed CPU/memory cap (see docker-compose.yml) so the only thing that can
// differ between formats is decode cost under identical resources.
//
// Measurement model: the topic already holds exactly -n messages (the constant
// lag set up by the producer). We start at offset 0, decode until we've seen n
// records, and report wall-clock time and messages/sec. Because every scenario
// drains the same fixed backlog under the same cap, RPS is directly comparable.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"serdebench/internal/codec"
	"serdebench/internal/event"
	"serdebench/internal/runutil"
)

func main() {
	brokers := flag.String("brokers", "localhost:9092", "comma-separated Kafka brokers")
	topic := flag.String("topic", "", "topic to drain (default: bench.<format>)")
	format := flag.String("format", "json", "serialization format: "+runutil.FmtList())
	n := flag.Int("n", 2_000_000, "number of messages to consume (must match the fixed lag)")
	flag.Parse()

	c, err := codec.Get(*format)
	if err != nil {
		log.Fatal(err)
	}
	tp := runutil.TopicName(*topic, *format)

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(runutil.SplitCSV(*brokers)...),
		kgo.ConsumeTopics(tp),
		// Always start from the beginning so the full fixed backlog is drained.
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
		kgo.FetchMaxBytes(50<<20),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cl.Close()

	ctx := context.Background()
	var (
		consumed int
		out      event.Event
		start    time.Time // started on first record so connection setup isn't counted
	)
	log.Printf("draining %d %s messages from %q ...", *n, *format, tp)
	for consumed < *n {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Fatalf("fetch: %v", errs)
		}
		fetches.EachRecord(func(r *kgo.Record) {
			if consumed == 0 {
				start = time.Now() // start timing on first record, excluding connect/setup
			}
			if err := c.Unmarshal(r.Value, &out); err != nil {
				log.Fatalf("unmarshal at %d: %v", consumed, err)
			}
			consumed++
		})
	}
	elapsed := time.Since(start)
	rps := float64(consumed) / elapsed.Seconds()
	log.Printf("RESULT format=%s n=%d elapsed=%s rps=%.0f ns/msg=%.0f",
		*format, consumed, elapsed.Round(time.Millisecond), rps,
		float64(elapsed.Nanoseconds())/float64(consumed))
}
