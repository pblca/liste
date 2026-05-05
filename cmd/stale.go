package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

var (
	staleDays int
)

var staleCmd = &cobra.Command{
	Use:   "stale",
	Short: "Show items that haven't been updated recently",
	Long: `Lists items that haven't been updated in N days (default 14).
Surfaces forgotten or abandoned work that may need attention.
Only shows items that are not done or cancelled.`,
	Args: cobra.NoArgs,
	RunE: runStale,
}

func init() {
	staleCmd.Flags().IntVar(&staleDays, "days", 14, "Days without update to consider stale")
	rootCmd.AddCommand(staleCmd)
}

func runStale(cmd *cobra.Command, args []string) error {
	allItems, _, _, err := collectAllItems()
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -staleDays)

	var stale []itemWithProject
	for _, entry := range allItems {
		item := entry.item

		// Skip done/cancelled
		if item.Status == "done" || item.Status == "cancelled" {
			continue
		}

		if item.Updated.Before(cutoff) {
			stale = append(stale, entry)
		}
	}

	// Sort by staleness (oldest update first)
	sort.SliceStable(stale, func(i, j int) bool {
		return stale[i].item.Updated.Before(stale[j].item.Updated)
	})

	if flagJSON {
		type jsonStale struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Title      string `json:"title"`
			Status     string `json:"status"`
			Priority   string `json:"priority"`
			Phase      *int   `json:"phase,omitempty"`
			Project    string `json:"project"`
			LastUpdate string `json:"last_update"`
			DaysStale  int    `json:"days_stale"`
		}
		var out []jsonStale
		for _, entry := range stale {
			days := int(time.Since(entry.item.Updated).Hours() / 24)
			out = append(out, jsonStale{
				ID:         entry.item.ID,
				Type:       string(entry.item.Type),
				Title:      entry.item.Title,
				Status:     entry.item.Status,
				Priority:   entry.item.Priority,
				Phase:      entry.item.Phase,
				Project:    entry.project,
				LastUpdate: entry.item.Updated.Format("2006-01-02"),
				DaysStale:  days,
			})
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(out)
		return nil
	}

	if len(stale) == 0 {
		fmt.Fprintf(os.Stdout, "No stale items (threshold: %d days)\n", staleDays)
		return nil
	}

	fmt.Fprintf(os.Stdout, "Stale items (not updated in %d+ days): %d\n\n", staleDays, len(stale))
	fmt.Fprintf(os.Stdout, "%-10s %-8s %-10s %-12s %-12s %s\n", "ID", "STATUS", "PRIORITY", "LAST UPDATE", "PROJECT", "TITLE")
	fmt.Fprintf(os.Stdout, "%-10s %-8s %-10s %-12s %-12s %s\n", "---", "------", "--------", "-----------", "-------", "-----")
	for _, entry := range stale {
		item := entry.item
		days := int(time.Since(item.Updated).Hours() / 24)
		fmt.Fprintf(os.Stdout, "%-10s %-8s %-10s %-12s %-12s %s (%dd ago)\n",
			item.ID, item.Status, item.Priority, item.Updated.Format("2006-01-02"), entry.project, item.Title, days)
	}

	return nil
}
