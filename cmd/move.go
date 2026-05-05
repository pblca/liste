package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move <id> <status>",
	Short: "Move an item to a new status",
	Long:  "Change the status of an item. Equivalent to 'liste set <id> status <status>'.",
	Args:  cobra.ExactArgs(2),
	RunE:  runMove,
}

func init() {
	rootCmd.AddCommand(moveCmd)
}

func runMove(cmd *cobra.Command, args []string) error {
	s, err := getStore()
	if err != nil {
		return err
	}

	id := strings.ToUpper(args[0])
	status := strings.ToLower(args[1])

	cfg, err := s.ReadConfig()
	if err != nil {
		return err
	}

	if !cfg.IsValidStatus(status) {
		return fmt.Errorf("invalid status %q (valid: %s)", status, strings.Join(cfg.Statuses, ", "))
	}

	item, err := s.ReadItem(id)
	if err != nil {
		return err
	}

	item.Status = status
	item.Updated = time.Now()

	if err := s.WriteItem(item); err != nil {
		return err
	}

	f := getFormatter()
	f.Message(fmt.Sprintf("Moved %s to %s: %s", id, status, item.Title))
	return nil
}
