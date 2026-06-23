// Package codec holds one implementation per serialization format. Every codec
// shares the same interface and the same input type (*event.Event), so the
// benchmark measures the format and nothing else.
//
// The formats are chosen to decompose the question "why is protobuf faster than
// JSON?" into two independent variables:
//
//	                     | binary? | keys on wire?
//	  json               |  text   |  yes (full names)
//	  jsonshort          |  text   |  yes (1-char names)   <- isolates key LENGTH
//	  msgpack            |  binary |  yes (full names)     <- the CONTROL: binary but keyed
//	  cbor               |  binary |  yes (full names)     <- second binary-keyed control
//	  protobuf           |  binary |  no  (field tags)
//	  avro               |  binary |  no  (positional)
//
// If "no keys" were the dominant factor, msgpack/cbor (binary WITH keys) would
// sit much closer to JSON than to protobuf/avro. If they sit close to
// protobuf/avro, then going binary — not dropping keys — is the real win.
package codec

import (
	"fmt"
	"sort"

	"serdebench/internal/event"
)

// Codec encodes and decodes a single Event. Implementations must be safe to
// reuse across calls (the benchmark holds one instance per format).
type Codec interface {
	Name() string
	Marshal(e *event.Event) ([]byte, error)
	Unmarshal(data []byte, e *event.Event) error
}

var registry = map[string]Codec{}

func register(c Codec) {
	if _, dup := registry[c.Name()]; dup {
		panic("duplicate codec: " + c.Name())
	}
	registry[c.Name()] = c
}

// Get returns the codec registered under name, or an error listing the
// available names.
func Get(name string) (Codec, error) {
	c, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown codec %q; available: %v", name, Names())
	}
	return c, nil
}

// All returns every registered codec, in stable name order.
func All() []Codec {
	out := make([]Codec, 0, len(registry))
	for _, c := range registry {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out
}

// Names returns every registered codec name, sorted.
func Names() []string {
	out := make([]string, 0, len(registry))
	for n := range registry {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
