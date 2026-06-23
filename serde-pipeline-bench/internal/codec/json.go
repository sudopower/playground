package codec

import (
	"encoding/json"

	"serdebench/internal/event"
)

func init() { register(jsonCodec{}) }

// jsonCodec is the baseline: text format, full field names on the wire, generic
// reflection-based encode/decode from the standard library.
type jsonCodec struct{}

func (jsonCodec) Name() string { return "json" }

func (jsonCodec) Marshal(e *event.Event) ([]byte, error) { return json.Marshal(e) }

func (jsonCodec) Unmarshal(data []byte, e *event.Event) error { return json.Unmarshal(data, e) }
