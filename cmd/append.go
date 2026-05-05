package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var appendCmd = &cobra.Command{
	Use:   "append <id> <text>",
	Short: "Append a note to an item's body",
	Long:  "Appends timestamped text to the item's markdown body.",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runAppend,
}

func init() {
	rootCmd.AddCommand(appendCmd)
}

func runAppend(cmd *cobra.Command, args []string) error {
	s, err := getStore()
	if err != nil {
		return err
	}

	id := strings.ToUpper(args[0])
	text := strings.Join(args[1:], " ")

	item, err := s.ReadItem(id)
	if err != nil {
		return err
	}

	// Append as a dated note
	note := fmt.Sprintf("- %s: %s\n", time.Now().Format("2006-01-02"), text)

	if item.Body == "" {
		item.Body = "## Notes\n\n" + note
	} else {
		item.Body += "\n" + note
	}

	item.Updated = time.Now()

	if err := s.WriteItem(item); err != nil {
		return err
	}

	f := getFormatter()
	f.Message(fmt.Sprintf("Appended note to %s", id))
	return nil
}
