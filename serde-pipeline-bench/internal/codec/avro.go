package codec

import (
	"github.com/hamba/avro/v2"

	"serdebench/internal/event"
)

func init() { register(newAvroCodec()) }

// avroCodec is binary AND keyless like protobuf, but more extreme: Avro carries
// no field tags at all. The schema fixes field order, so the payload is just the
// values back to back. If "drop the keys" were the whole story, avro and
// protobuf should be near-identical and both far ahead of the keyed formats.
//
// Avro's catch (and a good thing to call out in the post): the reader MUST have
// the schema. Here we hold it in process; in a real pipeline it comes from a
// schema registry, which adds its own lookup/caching cost.
type avroCodec struct {
	schema avro.Schema
}

const avroSchema = `{
  "type": "record",
  "name": "Event",
  "fields": [
    {"name": "device_id",    "type": "string"},
    {"name": "timestamp",    "type": "long"},
    {"name": "temperature",  "type": "double"},
    {"name": "humidity",     "type": "double"},
    {"name": "pressure",     "type": "double"},
    {"name": "latitude",     "type": "double"},
    {"name": "longitude",    "type": "double"},
    {"name": "battery_pct",  "type": "int"},
    {"name": "firmware",     "type": "string"},
    {"name": "status",       "type": "string"},
    {"name": "sequence_num", "type": "long"},
    {"name": "region",       "type": "string"}
  ]
}`

func newAvroCodec() avroCodec {
	return avroCodec{schema: avro.MustParse(avroSchema)}
}

func (avroCodec) Name() string { return "avro" }

func (c avroCodec) Marshal(e *event.Event) ([]byte, error) { return avro.Marshal(c.schema, e) }

func (c avroCodec) Unmarshal(data []byte, e *event.Event) error {
	return avro.Unmarshal(c.schema, data, e)
}
