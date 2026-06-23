package codec

import (
	"github.com/vmihailenco/msgpack/v5"

	"serdebench/internal/event"
)

func init() { register(msgpackCodec{}) }

// msgpackCodec is THE control: a binary format that still carries field names as
// map keys (msgpack v5 encodes structs as maps keyed by name by default). It
// shares "binary, schema-less, keyed" with CBOR. Comparing it against protobuf
// /avro isolates the value of dropping keys while holding binary-ness fixed.
type msgpackCodec struct{}

func (msgpackCodec) Name() string { return "msgpack" }

func (msgpackCodec) Marshal(e *event.Event) ([]byte, error) { return msgpack.Marshal(e) }

func (msgpackCodec) Unmarshal(data []byte, e *event.Event) error { return msgpack.Unmarshal(data, e) }
