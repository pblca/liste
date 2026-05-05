package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Show phase-level progress and overall completion",
	Long:  "Display progress bars for each phase and overall project completion percentage.",
	Args:  cobra.NoArgs,
	RunE:  runProgress,
}

func init() {
	rootCmd.AddCommand(progressCmd)
}

type phaseProgress struct {
	Phase   int
	Done    int
	Total   int
	Active  int
	Blocked int
}

func runProgress(cmd *cobra.Command, args []string) error {
	allItems, _, projectName, err := collectAllItems()
	if err != nil {
		return err
	}

	// Compute per-phase progress
	phaseMap := make(map[int]*phaseProgress)
	var unphasedDone, unphasedTotal int
	totalDone := 0
	totalAll := 0

	for _, entry := range allItems {
		item := entry.item
		totalAll++

		if item.Phase == nil {
			unphasedTotal++
			if item.Status == "done" || item.Status == "cancelled" {
				unphasedDone++
				totalDone++
			}
			continue
		}

		phase := *item.Phase
		if phaseMap[phase] == nil {
			phaseMap[phase] = &phaseProgress{Phase: phase}
		}
		pp := phaseMap[phase]
		pp.Total++

		switch {
		case item.Status == "done" || item.Status == "cancelled":
			pp.Done++
			totalDone++
		case item.Status == "active":
			pp.Active++
		case item.Blocked != nil:
			pp.Blocked++
		}
	}

	// Sort phases
	var phases []int
	for p := range phaseMap {
		phases = append(phases, p)
	}
	sort.Ints(phases)

	if flagJSON {
		type jsonPhaseProgress struct {
			Phase   int     `json:"phase"`
			Done    int     `json:"done"`
			Total   int     `json:"total"`
			Active  int     `json:"active"`
			Blocked int     `json:"blocked"`
			Percent float64 `json:"percent"`
		}
		var jPhases []jsonPhaseProgress
		for _, p := range phases {
			pp := phaseMap[p]
			pct := 0.0
			if pp.Total > 0 {
				pct = float64(pp.Done) / float64(pp.Total) * 100
			}
			jPhases = append(jPhases, jsonPhaseProgress{
				Phase:   p,
				Done:    pp.Done,
				Total:   pp.Total,
				Active:  pp.Active,
				Blocked: pp.Blocked,
				Percent: pct,
			})
		}
		overallPct := 0.0
		if totalAll > 0 {
			overallPct = float64(totalDone) / float64(totalAll) * 100
		}
		out := map[string]any{
			"project":         projectName,
			"total_items":     totalAll,
			"total_done":      totalDone,
			"overall_percent": overallPct,
			"phases":          jPhases,
			"unphased_done":   unphasedDone,
			"unphased_total":  unphasedTotal,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(out)
		return nil
	}

	// Text output with progress bars
	overallPct := 0.0
	if totalAll > 0 {
		overallPct = float64(totalDone) / float64(totalAll) * 100
	}

	fmt.Fprintf(os.Stdout, "Project: %s\n", projectName)
	fmt.Fprintf(os.Stdout, "Overall: %d/%d (%.0f%%)\n", totalDone, totalAll, overallPct)
	fmt.Fprintf(os.Stdout, "%s\n\n", progressBar(totalDone, totalAll, 40))

	for _, p := range phases {
		pp := phaseMap[p]
		pct := 0.0
		if pp.Total > 0 {
			pct = float64(pp.Done) / float64(pp.Total) * 100
		}

		status := ""
		if pp.Done == pp.Total {
			status = " (complete)"
		} else if pp.Active > 0 {
			status = " (active)"
		}

		fmt.Fprintf(os.Stdout, "Phase %d: %d/%d (%.0f%%)%s\n", p, pp.Done, pp.Total, pct, status)
		fmt.Fprintf(os.Stdout, "  %s", progressBar(pp.Done, pp.Total, 30))
		if pp.Active > 0 || pp.Blocked > 0 {
			details := []string{}
			if pp.Active > 0 {
				details = append(details, fmt.Sprintf("%d active", pp.Active))
			}
			if pp.Blocked > 0 {
				details = append(details, fmt.Sprintf("%d blocked", pp.Blocked))
			}
			fmt.Fprintf(os.Stdout, "  [%s]", strings.Join(details, ", "))
		}
		fmt.Fprintln(os.Stdout)
	}

	if unphasedTotal > 0 {
		fmt.Fprintf(os.Stdout, "\nUnphased: %d/%d\n", unphasedDone, unphasedTotal)
	}

	return nil
}

// progressBar renders an ASCII progress bar.
func progressBar(done, total, width int) string {
	if total == 0 {
		return "[" + strings.Repeat(" ", width) + "]"
	}

	filled := (done * width) / total
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return "[" + bar + "]"
}
