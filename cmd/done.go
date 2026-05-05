package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Mark an item as done",
	Long:  "Shortcut to set status to 'done' and clear any blocked flag.",
	Args:  cobra.ExactArgs(1),
	RunE:  runDone,
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

func runDone(cmd *cobra.Command, args []string) error {
	s, err := getStore()
	if err != nil {
		return err
	}

	id := strings.ToUpper(args[0])
	item, err := s.ReadItem(id)
	if err != nil {
		return err
	}

	item.Status = "done"
	item.Blocked = nil
	item.Updated = time.Now()

	if err := s.WriteItem(item); err != nil {
		return err
	}

	f := getFormatter()
	f.Message(fmt.Sprintf("Marked %s as done: %s", id, item.Title))
	return nil
}
