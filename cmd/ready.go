package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var readyCmd = &cobra.Command{
	Use:   "ready",
	Short: "List all items whose dependencies are satisfied",
	Long: `Shows the "ready queue" — items that are not done, not blocked, and whose
depends-on targets are all complete. These are items that can be worked on now.`,
	Args: cobra.NoArgs,
	RunE: runReady,
}

func init() {
	rootCmd.AddCommand(readyCmd)
}

func runReady(cmd *cobra.Command, args []string) error {
	allItems, allItemsByID, _, err := collectAllItems()
	if err != nil {
		return err
	}

	var ready []itemWithProject
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

		ready = append(ready, entry)
	}

	// Sort by phase then priority
	sort.SliceStable(ready, func(i, j int) bool {
		a := ready[i].item
		b := ready[j].item
		if phaseOrder(a) != phaseOrder(b) {
			return phaseOrder(a) < phaseOrder(b)
		}
		return priorityWeight(a.Priority) < priorityWeight(b.Priority)
	})

	if flagJSON {
		type jsonItem struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			Priority string `json:"priority"`
			Phase    *int   `json:"phase,omitempty"`
			Project  string `json:"project"`
		}
		var out []jsonItem
		for _, entry := range ready {
			out = append(out, jsonItem{
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
		_ = enc.Encode(out)
		return nil
	}

	if flagQuiet {
		for _, entry := range ready {
			fmt.Fprintln(os.Stdout, entry.item.ID)
		}
		return nil
	}

	if len(ready) == 0 {
		fmt.Fprintln(os.Stdout, "No items ready to work on.")
		return nil
	}

	fmt.Fprintf(os.Stdout, "Ready queue (%d items):\n\n", len(ready))
	fmt.Fprintf(os.Stdout, "%-10s %-8s %-8s %-10s %-12s %s\n", "ID", "TYPE", "PHASE", "PRIORITY", "PROJECT", "TITLE")
	fmt.Fprintf(os.Stdout, "%-10s %-8s %-8s %-10s %-12s %s\n", "---", "----", "-----", "--------", "-------", "-----")
	for _, entry := range ready {
		item := entry.item
		phaseStr := "-"
		if item.Phase != nil {
			phaseStr = fmt.Sprintf("%d", *item.Phase)
		}
		fmt.Fprintf(os.Stdout, "%-10s %-8s %-8s %-10s %-12s %s\n",
			item.ID, item.Type, phaseStr, item.Priority, entry.project, item.Title)
	}

	return nil
}
