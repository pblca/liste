package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <id> <field> <value>",
	Short: "Set a field value on an item",
	Long:  "Set status, priority, title, tags, or phase on an item. Field names: status, priority, title, tags, phase.",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runSet,
}

func init() {
	rootCmd.AddCommand(setCmd)
}

func runSet(cmd *cobra.Command, args []string) error {
	s, err := getStore()
	if err != nil {
		return err
	}

	id := strings.ToUpper(args[0])
	field := strings.ToLower(args[1])
	value := ""
	if len(args) > 2 {
		value = strings.Join(args[2:], " ")
	}

	item, err := s.ReadItem(id)
	if err != nil {
		return err
	}

	cfg, err := s.ReadConfig()
	if err != nil {
		return err
	}

	switch field {
	case "status":
		if !cfg.IsValidStatus(value) {
			return fmt.Errorf("invalid status %q (valid: %s)", value, strings.Join(cfg.Statuses, ", "))
		}
		item.Status = value
	case "priority":
		if !cfg.IsValidPriority(value) {
			return fmt.Errorf("invalid priority %q (valid: %s)", value, strings.Join(cfg.Priorities, ", "))
		}
		item.Priority = value
	case "title":
		if value == "" {
			return fmt.Errorf("title cannot be empty")
		}
		item.Title = value
	case "tags":
		item.Tags = strings.Split(value, ",")
		for i := range item.Tags {
			item.Tags[i] = strings.TrimSpace(item.Tags[i])
		}
	case "phase":
		if value == "" || value == "0" {
			item.Phase = nil
		} else {
			p, err := strconv.Atoi(value)
			if err != nil || p < 1 {
				return fmt.Errorf("phase must be a positive integer (or empty to clear)")
			}
			item.Phase = &p
		}
	default:
		return fmt.Errorf("unknown field %q (valid: status, priority, title, tags, phase)", field)
	}

	item.Updated = time.Now()
	if err := s.WriteItem(item); err != nil {
		return err
	}

	f := getFormatter()
	f.Message(fmt.Sprintf("Updated %s: %s = %s", id, field, value))
	return nil
}
