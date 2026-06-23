package codec

import (
	"fmt"
	"math"

	"google.golang.org/protobuf/encoding/protowire"

	"serdebench/internal/event"
)

func init() { register(protobufCodec{}) }

// protobufCodec encodes the Event exactly as protobuf does on the wire, using
// the official low-level protowire package. We hand-roll the encode/decode for
// our fixed schema instead of generating code with protoc — the bytes are
// identical to what generated code emits (generated code calls protowire too),
// and it keeps the repo buildable with plain `go build`, no protoc toolchain.
//
// Schema (field numbers are the "keys", encoded as compact varint tags):
//
//	1: device_id string   5: pressure  double   9: firmware string
//	2: timestamp int64     6: latitude  double  10: status   string
//	3: temperature double  7: longitude double  11: sequence_num int64
//	4: humidity double     8: battery_pct int32 12: region   string
type protobufCodec struct{}

func (protobufCodec) Name() string { return "protobuf" }

func (protobufCodec) Marshal(e *event.Event) ([]byte, error) {
	var b []byte
	b = appendString(b, 1, e.DeviceID)
	b = appendVarint(b, 2, uint64(e.Timestamp))
	b = appendDouble(b, 3, e.Temperature)
	b = appendDouble(b, 4, e.Humidity)
	b = appendDouble(b, 5, e.Pressure)
	b = appendDouble(b, 6, e.Latitude)
	b = appendDouble(b, 7, e.Longitude)
	b = appendVarint(b, 8, uint64(uint32(e.BatteryPct)))
	b = appendString(b, 9, e.Firmware)
	b = appendString(b, 10, e.Status)
	b = appendVarint(b, 11, uint64(e.SequenceNum))
	b = appendString(b, 12, e.Region)
	return b, nil
}

func (protobufCodec) Unmarshal(data []byte, e *event.Event) error {
	*e = event.Event{}
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return fmt.Errorf("protobuf: bad tag: %w", protowire.ParseError(n))
		}
		data = data[n:]

		switch typ {
		case protowire.VarintType:
			v, m := protowire.ConsumeVarint(data)
			if m < 0 {
				return fmt.Errorf("protobuf: bad varint: %w", protowire.ParseError(m))
			}
			data = data[m:]
			switch num {
			case 2:
				e.Timestamp = int64(v)
			case 8:
				e.BatteryPct = int32(uint32(v))
			case 11:
				e.SequenceNum = int64(v)
			}
		case protowire.Fixed64Type:
			v, m := protowire.ConsumeFixed64(data)
			if m < 0 {
				return fmt.Errorf("protobuf: bad fixed64: %w", protowire.ParseError(m))
			}
			data = data[m:]
			f := math.Float64frombits(v)
			switch num {
			case 3:
				e.Temperature = f
			case 4:
				e.Humidity = f
			case 5:
				e.Pressure = f
			case 6:
				e.Latitude = f
			case 7:
				e.Longitude = f
			}
		case protowire.BytesType:
			v, m := protowire.ConsumeBytes(data)
			if m < 0 {
				return fmt.Errorf("protobuf: bad bytes: %w", protowire.ParseError(m))
			}
			data = data[m:]
			s := string(v)
			switch num {
			case 1:
				e.DeviceID = s
			case 9:
				e.Firmware = s
			case 10:
				e.Status = s
			case 12:
				e.Region = s
			}
		default:
			m := protowire.ConsumeFieldValue(num, typ, data)
			if m < 0 {
				return fmt.Errorf("protobuf: bad field: %w", protowire.ParseError(m))
			}
			data = data[m:]
		}
	}
	return nil
}

func appendString(b []byte, num protowire.Number, s string) []byte {
	if s == "" {
		return b
	}
	b = protowire.AppendTag(b, num, protowire.BytesType)
	return protowire.AppendString(b, s)
}

func appendVarint(b []byte, num protowire.Number, v uint64) []byte {
	if v == 0 {
		return b
	}
	b = protowire.AppendTag(b, num, protowire.VarintType)
	return protowire.AppendVarint(b, v)
}

func appendDouble(b []byte, num protowire.Number, f float64) []byte {
	b = protowire.AppendTag(b, num, protowire.Fixed64Type)
	return protowire.AppendFixed64(b, math.Float64bits(f))
}
