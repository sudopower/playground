// Package event defines the single canonical record used across every codec and
// the pipeline. Keeping ONE struct shared by all formats is what makes the
// benchmark fair: the only thing that varies between scenarios is the
// serialization, never the data model.
package event

import (
	"math/rand"
)

// Event is a representative telemetry record: a mix of strings, floats and
// integers, with realistic (i.e. not artificially short) field names. The field
// names matter — they are exactly the bytes that JSON carries on the wire and
// that protobuf/avro omit. That contrast is the whole point of the experiment.
type Event struct {
	DeviceID    string  `json:"device_id"    msgpack:"device_id"    cbor:"device_id"    avro:"device_id"`
	Timestamp   int64   `json:"timestamp"    msgpack:"timestamp"    cbor:"timestamp"    avro:"timestamp"`
	Temperature float64 `json:"temperature"  msgpack:"temperature"  cbor:"temperature"  avro:"temperature"`
	Humidity    float64 `json:"humidity"     msgpack:"humidity"     cbor:"humidity"     avro:"humidity"`
	Pressure    float64 `json:"pressure"     msgpack:"pressure"     cbor:"pressure"     avro:"pressure"`
	Latitude    float64 `json:"latitude"     msgpack:"latitude"     cbor:"latitude"     avro:"latitude"`
	Longitude   float64 `json:"longitude"    msgpack:"longitude"    cbor:"longitude"    avro:"longitude"`
	BatteryPct  int32   `json:"battery_pct"  msgpack:"battery_pct"  cbor:"battery_pct"  avro:"battery_pct"`
	Firmware    string  `json:"firmware"     msgpack:"firmware"     cbor:"firmware"     avro:"firmware"`
	Status      string  `json:"status"       msgpack:"status"       cbor:"status"       avro:"status"`
	SequenceNum int64   `json:"sequence_num" msgpack:"sequence_num" cbor:"sequence_num" avro:"sequence_num"`
	Region      string  `json:"region"       msgpack:"region"       cbor:"region"       avro:"region"`
}

var (
	firmwares = []string{"v2.3.1", "v2.4.0", "v3.0.0-rc1", "v1.9.7"}
	statuses  = []string{"ok", "degraded", "warning", "critical"}
	regions   = []string{"eu-central-1", "us-east-1", "ap-southeast-2", "us-west-2"}
)

// Sample builds a deterministic record from a seed so every format encodes the
// exact same logical data and runs are reproducible across scenarios.
func Sample(seed int64) Event {
	r := rand.New(rand.NewSource(seed))
	return Event{
		DeviceID:    "device-" + pad(int(seed%100000)),
		Timestamp:   1_700_000_000_000_000_000 + seed*1_000_000,
		Temperature: 15 + r.Float64()*20,
		Humidity:    30 + r.Float64()*50,
		Pressure:    980 + r.Float64()*60,
		Latitude:    -90 + r.Float64()*180,
		Longitude:   -180 + r.Float64()*360,
		BatteryPct:  int32(r.Intn(101)),
		Firmware:    firmwares[r.Intn(len(firmwares))],
		Status:      statuses[r.Intn(len(statuses))],
		SequenceNum: seed,
		Region:      regions[r.Intn(len(regions))],
	}
}

// Batch returns n deterministic samples.
func Batch(n int) []Event {
	out := make([]Event, n)
	for i := range out {
		out[i] = Sample(int64(i))
	}
	return out
}

func pad(n int) string {
	const width = 5
	s := []byte("00000")
	i := len(s) - 1
	for n > 0 && i >= 0 {
		s[i] = byte('0' + n%10)
		n /= 10
		i--
	}
	return string(s[len(s)-width:])
}
