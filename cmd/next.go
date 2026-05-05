package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var (
	nextCount int
)

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show the next item(s) ready to be worked on",
	Long: `Determines what to work on next based on:
1. Dependencies satisfied (all depends-on targets are done)
2. Not blocked
3. Earliest active phase first
4. Highest priority within phase
5. Status is not already active/done/cancelled

Use --count to show multiple candidates.`,
	Args: cobra.NoArgs,
	RunE: runNext,
}

func init() {
	nextCmd.Flags().IntVarP(&nextCount, "count", "n", 1, "Number of candidates to show")
	rootCmd.AddCommand(nextCmd)
}

func runNext(cmd *cobra.Command, args []string) error {
	allItems, allItemsByID, _, err := collectAllItems()
	if err != nil {
		return err
	}

	// Filter to candidates: not done, not cancelled, not blocked, not already active
	var candidates []itemWithProject
	for _, entry := range allItems {
		item := entry.item

		if item.Status == "done" || item.Status == "cancelled" || item.Status == "active" {
			continue
		}
		if item.Blocked != nil {
			continue
		}
		if !depsResolved(item, allItemsByID) {
			continue
		}

		candidates = append(candidates, entry)
	}

	if len(candidates) == 0 {
		f := getFormatter()
		f.Message("No items ready to work on.")
		return nil
	}

	// Sort candidates: phase (lowest first, nil last), then priority, then created date
	sort.SliceStable(candidates, func(i, j int) bool {
		a := candidates[i].item
		b := candidates[j].item

		aPhase := phaseOrder(a)
		bPhase := phaseOrder(b)
		if aPhase != bPhase {
			return aPhase < bPhase
		}

		aPri := priorityWeight(a.Priority)
		bPri := priorityWeight(b.Priority)
		if aPri != bPri {
			return aPri < bPri
		}

		return a.Created.Before(b.Created)
	})

	if nextCount > len(candidates) {
		nextCount = len(candidates)
	}
	selected := candidates[:nextCount]

	if flagJSON {
		renderNextJSON(selected)
		return nil
	}

	if flagQuiet {
		for _, entry := range selected {
			fmt.Fprintln(os.Stdout, entry.item.ID)
		}
		return nil
	}

	if len(selected) == 1 {
		entry := selected[0]
		item := entry.item
		phaseStr := ""
		if item.Phase != nil {
			phaseStr = fmt.Sprintf(" (phase %d)", *item.Phase)
		}
		fmt.Fprintf(os.Stdout, "Next: %s [%s] %s%s\n", item.ID, item.Priority, item.Title, phaseStr)
		fmt.Fprintf(os.Stdout, "  Type: %s | Status: %s | Project: %s\n", item.Type, item.Status, entry.project)
	} else {
		fmt.Fprintf(os.Stdout, "Next %d candidates:\n\n", len(selected))
		for i, entry := range selected {
			item := entry.item
			phaseStr := ""
			if item.Phase != nil {
				phaseStr = fmt.Sprintf(" p%d", *item.Phase)
			}
			fmt.Fprintf(os.Stdout, "  %d. %-10s [%-8s] %-8s %s%s\n",
				i+1, item.ID, item.Priority, entry.project, item.Title, phaseStr)
		}
	}

	return nil
}

func renderNextJSON(items []itemWithProject) {
	type jsonCandidate struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Title    string `json:"title"`
		Status   string `json:"status"`
		Priority string `json:"priority"`
		Phase    *int   `json:"phase,omitempty"`
		Project  string `json:"project"`
	}

	var candidates []jsonCandidate
	for _, entry := range items {
		candidates = append(candidates, jsonCandidate{
			ID:       entry.item.ID,
			Type:     string(entry.item.Type),
			Title:    entry.item.Title,
			Status:   entry.item.Status,
			Priority: entry.item.Priority,
			Phase:    entry.item.Phase,
			Project:  entry.project,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(candidates)
}
