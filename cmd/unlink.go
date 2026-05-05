package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink <id> <target>",
	Short: "Remove a link between items",
	Long:  "Remove all links from the source item to the target.",
	Args:  cobra.ExactArgs(2),
	RunE:  runUnlink,
}

func init() {
	rootCmd.AddCommand(unlinkCmd)
}

func runUnlink(cmd *cobra.Command, args []string) error {
	s, err := getStore()
	if err != nil {
		return err
	}

	id := strings.ToUpper(args[0])
	target := strings.ToUpper(args[1])

	item, err := s.ReadItem(id)
	if err != nil {
		return err
	}

	// Remove all links to target
	found := false
	newLinks := item.Links[:0]
	for _, link := range item.Links {
		if link.Target == target {
			found = true
			continue
		}
		newLinks = append(newLinks, link)
	}

	if !found {
		return fmt.Errorf("no link from %s to %s", id, target)
	}

	item.Links = newLinks
	item.Updated = time.Now()

	if err := s.WriteItem(item); err != nil {
		return err
	}

	f := getFormatter()
	f.Message(fmt.Sprintf("Removed link(s) from %s to %s", id, target))
	return nil
}
