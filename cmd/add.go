package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/pblca/liste/internal/model"
	"github.com/spf13/cobra"
)

var (
	addPriority string
	addTags     []string
	addStatus   string
	addPhase    int
)

var addCmd = &cobra.Command{
	Use:   "add <type> <title>",
	Short: "Create a new item",
	Long:  "Create a new item of the given type (feature, bug, idea, task, epic).",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runAdd,
}

func init() {
	addCmd.Flags().StringVar(&addPriority, "priority", "", "Priority (critical, high, medium, low)")
	addCmd.Flags().StringSliceVar(&addTags, "tag", nil, "Tags (can be specified multiple times)")
	addCmd.Flags().StringVar(&addStatus, "status", "", "Initial status (overrides default)")
	addCmd.Flags().IntVar(&addPhase, "phase", 0, "Phase number (0 = unphased)")
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	s, err := getStore()
	if err != nil {
		return err
	}

	typeStr := args[0]
	title := strings.Join(args[1:], " ")

	itemType, ok := model.ParseItemType(typeStr)
	if !ok {
		return fmt.Errorf("invalid type %q (valid: feature, bug, idea, task, epic)", typeStr)
	}

	cfg, err := s.ReadConfig()
	if err != nil {
		return err
	}

	item, err := s.CreateItem(itemType, title, cfg)
	if err != nil {
		return err
	}

	// Apply optional overrides
	changed := false
	if addPriority != "" {
		if !cfg.IsValidPriority(addPriority) {
			return fmt.Errorf("invalid priority %q (valid: %s)", addPriority, strings.Join(cfg.Priorities, ", "))
		}
		item.Priority = addPriority
		changed = true
	}
	if addStatus != "" {
		if !cfg.IsValidStatus(addStatus) {
			return fmt.Errorf("invalid status %q (valid: %s)", addStatus, strings.Join(cfg.Statuses, ", "))
		}
		item.Status = addStatus
		changed = true
	}
	if len(addTags) > 0 {
		item.Tags = addTags
		changed = true
	}
	if addPhase > 0 {
		p := addPhase
		item.Phase = &p
		changed = true
	}
	if changed {
		item.Updated = time.Now()
		if err := s.WriteItem(item); err != nil {
			return err
		}
	}

	f := getFormatter()
	f.ItemCreated(item)
	return nil
}
