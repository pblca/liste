package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	diffSince string
	diffDays  int
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show what changed since a given date",
	Long: `Shows items that were created or updated since a given date.
Useful for standup summaries or AI agent "catch-up".

Use --since with a date (YYYY-MM-DD) or --days for relative lookback.`,
	Args: cobra.NoArgs,
	RunE: runDiff,
}

func init() {
	diffCmd.Flags().StringVar(&diffSince, "since", "", "Show changes since date (YYYY-MM-DD)")
	diffCmd.Flags().IntVar(&diffDays, "days", 1, "Show changes from last N days (default: 1)")
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	allItems, _, _, err := collectAllItems()
	if err != nil {
		return err
	}

	// Determine cutoff time
	var since time.Time
	if diffSince != "" {
		parsed, err := time.Parse("2006-01-02", diffSince)
		if err != nil {
			return fmt.Errorf("invalid date %q (use YYYY-MM-DD)", diffSince)
		}
		since = parsed
	} else {
		since = time.Now().AddDate(0, 0, -diffDays)
	}

	// Truncate to start of day
	since = time.Date(since.Year(), since.Month(), since.Day(), 0, 0, 0, 0, since.Location())

	var created []itemWithProject
	var updated []itemWithProject
	var completed []itemWithProject

	for _, entry := range allItems {
		item := entry.item

		itemCreated := time.Date(item.Created.Year(), item.Created.Month(), item.Created.Day(), 0, 0, 0, 0, item.Created.Location())
		itemUpdated := time.Date(item.Updated.Year(), item.Updated.Month(), item.Updated.Day(), 0, 0, 0, 0, item.Updated.Location())

		if !itemCreated.Before(since) {
			created = append(created, entry)
		} else if !itemUpdated.Before(since) {
			if item.Status == "done" {
				completed = append(completed, entry)
			} else {
				updated = append(updated, entry)
			}
		}
	}

	if flagJSON {
		type jsonDiffItem struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Title    string `json:"title"`
			Status   string `json:"status"`
			Priority string `json:"priority"`
			Phase    *int   `json:"phase,omitempty"`
			Project  string `json:"project"`
		}
		toJSON := func(items []itemWithProject) []jsonDiffItem {
			var out []jsonDiffItem
			for _, entry := range items {
				out = append(out, jsonDiffItem{
					ID:       entry.item.ID,
					Type:     string(entry.item.Type),
					Title:    entry.item.Title,
					Status:   entry.item.Status,
					Priority: entry.item.Priority,
					Phase:    entry.item.Phase,
					Project:  entry.project,
				})
			}
			return out
		}
		out := map[string]any{
			"since":     since.Format("2006-01-02"),
			"created":   toJSON(created),
			"updated":   toJSON(updated),
			"completed": toJSON(completed),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(out)
		return nil
	}

	totalChanges := len(created) + len(updated) + len(completed)
	if totalChanges == 0 {
		fmt.Fprintf(os.Stdout, "No changes since %s\n", since.Format("2006-01-02"))
		return nil
	}

	fmt.Fprintf(os.Stdout, "Changes since %s:\n\n", since.Format("2006-01-02"))

	if len(completed) > 0 {
		fmt.Fprintf(os.Stdout, "COMPLETED (%d)\n", len(completed))
		for _, entry := range completed {
			fmt.Fprintf(os.Stdout, "  %s %s (%s)\n", entry.item.ID, entry.item.Title, entry.project)
		}
		fmt.Fprintln(os.Stdout)
	}

	if len(created) > 0 {
		fmt.Fprintf(os.Stdout, "CREATED (%d)\n", len(created))
		for _, entry := range created {
			fmt.Fprintf(os.Stdout, "  %s [%s] %s (%s)\n", entry.item.ID, entry.item.Priority, entry.item.Title, entry.project)
		}
		fmt.Fprintln(os.Stdout)
	}

	if len(updated) > 0 {
		fmt.Fprintf(os.Stdout, "UPDATED (%d)\n", len(updated))
		for _, entry := range updated {
			fmt.Fprintf(os.Stdout, "  %s [%s] %s (%s)\n", entry.item.ID, entry.item.Status, entry.item.Title, entry.project)
		}
		fmt.Fprintln(os.Stdout)
	}

	return nil
}
