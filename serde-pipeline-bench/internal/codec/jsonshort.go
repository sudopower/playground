package codec

import (
	"encoding/json"

	"serdebench/internal/event"
)

func init() { register(jsonShortCodec{}) }

// jsonShortCodec is JSON with single-character keys. It holds everything about
// the baseline constant (still text, still reflection) EXCEPT key length, so the
// json -> jsonshort delta isolates how much of JSON's cost is just carrying long
// field names versus the cost of text parsing itself.
type jsonShortCodec struct{}

func (jsonShortCodec) Name() string { return "jsonshort" }

// shortEvent mirrors event.Event field-for-field with terse keys.
type shortEvent struct {
	D  string  `json:"a"`
	Ts int64   `json:"b"`
	Te float64 `json:"c"`
	H  float64 `json:"d"`
	P  float64 `json:"e"`
	La float64 `json:"f"`
	Lo float64 `json:"g"`
	B  int32   `json:"h"`
	Fw string  `json:"i"`
	St string  `json:"j"`
	Sq int64   `json:"k"`
	R  string  `json:"l"`
}

func (jsonShortCodec) Marshal(e *event.Event) ([]byte, error) {
	return json.Marshal(shortEvent{
		D: e.DeviceID, Ts: e.Timestamp, Te: e.Temperature, H: e.Humidity,
		P: e.Pressure, La: e.Latitude, Lo: e.Longitude, B: e.BatteryPct,
		Fw: e.Firmware, St: e.Status, Sq: e.SequenceNum, R: e.Region,
	})
}

func (jsonShortCodec) Unmarshal(data []byte, e *event.Event) error {
	var s shortEvent
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	e.DeviceID, e.Timestamp, e.Temperature, e.Humidity = s.D, s.Ts, s.Te, s.H
	e.Pressure, e.Latitude, e.Longitude, e.BatteryPct = s.P, s.La, s.Lo, s.B
	e.Firmware, e.Status, e.SequenceNum, e.Region = s.Fw, s.St, s.Sq, s.R
	return nil
}
