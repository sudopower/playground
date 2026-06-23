package codec

import (
	"testing"

	"serdebench/internal/event"
)

// TestRoundTrip guards the experiment's integrity: every codec must reproduce
// the record exactly. A benchmark of a codec that silently drops a field would
// be meaningless, so this runs before we trust any number.
func TestRoundTrip(t *testing.T) {
	want := event.Sample(42)
	for _, c := range All() {
		c := c
		t.Run(c.Name(), func(t *testing.T) {
			data, err := c.Marshal(&want)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var got event.Event
			if err := c.Unmarshal(data, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got != want {
				t.Fatalf("round-trip mismatch\n want=%+v\n  got=%+v", want, got)
			}
		})
	}
}
