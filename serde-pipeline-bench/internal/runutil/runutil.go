// Package runutil holds small helpers shared by the producer and consumer
// commands (flag parsing, topic naming) so the two stay consistent.
package runutil

import (
	"strings"

	"serdebench/internal/codec"
)

// SplitCSV turns "a,b,c" into []string{"a","b","c"}, trimming spaces.
func SplitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// FmtList is the human-readable list of registered formats, for flag help.
func FmtList() string { return strings.Join(codec.Names(), "|") }

// TopicName resolves the topic: an explicit override, else bench.<format> so
// each format drains its own dedicated, isolated topic.
func TopicName(override, format string) string {
	if override != "" {
		return override
	}
	return "bench." + format
}
