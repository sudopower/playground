package codec

import (
	"github.com/fxamacker/cbor/v2"

	"serdebench/internal/event"
)

func init() { register(cborCodec{}) }

// cborCodec is a second binary-but-keyed control. By default fxamacker/cbor
// encodes structs as maps keyed by field name, so like msgpack it pays for keys
// on the wire but skips text parsing. Two independent keyed-binary formats guard
// against a result that's really just one library's implementation quirk.
type cborCodec struct{}

func (cborCodec) Name() string { return "cbor" }

func (cborCodec) Marshal(e *event.Event) ([]byte, error) { return cbor.Marshal(e) }

func (cborCodec) Unmarshal(data []byte, e *event.Event) error { return cbor.Unmarshal(data, e) }
