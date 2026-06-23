// Command sizes reports the encoded wire size of one representative record for
// every format. Size is the variable the "no keys" hypothesis is really about,
// so we measure it directly and separately from speed.
package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"serdebench/internal/codec"
	"serdebench/internal/event"
)

func main() {
	// Average over a batch so variable-length strings don't skew a single sample.
	const n = 1000
	batch := event.Batch(n)

	type row struct {
		name  string
		total int
	}
	var rows []row
	for _, c := range codec.All() {
		total := 0
		for i := range batch {
			data, err := c.Marshal(&batch[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", c.Name(), err)
				os.Exit(1)
			}
			total += len(data)
		}
		rows = append(rows, row{c.Name(), total})
	}

	// Use json as the 100% baseline for the relative column.
	var base float64
	for _, r := range rows {
		if r.name == "json" {
			base = float64(r.total)
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "FORMAT\tBYTES/MSG (avg)\tVS JSON")
	for _, r := range rows {
		avg := float64(r.total) / float64(n)
		pct := 100 * float64(r.total) / base
		fmt.Fprintf(w, "%s\t%.1f\t%.0f%%\n", r.name, avg, pct)
	}
	w.Flush()
}
