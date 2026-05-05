package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an item",
	Long:  "Permanently remove an item file. Use --force to skip confirmation.",
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

var deleteForce bool

func init() {
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation")
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	s, err := getStore()
	if err != nil {
		return err
	}

	id := strings.ToUpper(args[0])
	item, err := s.ReadItem(id)
	if err != nil {
		return err
	}

	if !deleteForce {
		// In non-interactive mode (AI usage), require --force
		fmt.Printf("Delete %s (%s)? Use --force to confirm.\n", id, item.Title)
		return nil
	}

	if err := s.DeleteItem(id); err != nil {
		return err
	}

	f := getFormatter()
	f.Message(fmt.Sprintf("Deleted %s: %s", id, item.Title))
	return nil
}
