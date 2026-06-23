// Command producer fills a Kafka topic with exactly -n messages encoded in a
// chosen format, then exits. It does NOT run concurrently with the consumer.
//
// This is how we hold "Kafka lag" constant across scenarios: every scenario
// starts with the same fixed backlog (lag == n, consumer offset 0) sitting in
// the topic. We then start the resource-capped consumer and time how long it
// takes to drain to zero. Same starting lag, same record count, same data — only
// the wire format differs.
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
	topic := flag.String("topic", "", "topic to fill (default: bench.<format>)")
	format := flag.String("format", "json", "serialization format: "+runutil.FmtList())
	n := flag.Int("n", 2_000_000, "number of messages to produce (the fixed lag)")
	flag.Parse()

	c, err := codec.Get(*format)
	if err != nil {
		log.Fatal(err)
	}
	tp := runutil.TopicName(*topic, *format)

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(runutil.SplitCSV(*brokers)...),
		kgo.DefaultProduceTopic(tp),
		kgo.ProducerBatchMaxBytes(16<<20),
		kgo.MaxBufferedRecords(250_000),
		kgo.RequiredAcks(kgo.LeaderAck()),
		kgo.ProducerLinger(50*time.Millisecond),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cl.Close()

	ctx := context.Background()
	start := time.Now()
	var produced int64
	for i := 0; i < *n; i++ {
		e := event.Sample(int64(i))
		val, err := c.Marshal(&e)
		if err != nil {
			log.Fatalf("marshal: %v", err)
		}
		cl.Produce(ctx, &kgo.Record{Key: []byte(e.DeviceID), Value: val}, func(_ *kgo.Record, err error) {
			if err != nil {
				log.Fatalf("produce: %v", err)
			}
		})
		produced++
		if produced%200_000 == 0 {
			log.Printf("produced %d/%d", produced, *n)
		}
	}
	if err := cl.Flush(ctx); err != nil {
		log.Fatalf("flush: %v", err)
	}
	log.Printf("done: %d %s messages to %q in %s (lag is now %d)",
		produced, *format, tp, time.Since(start).Round(time.Millisecond), produced)
}
