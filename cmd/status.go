package cmd

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show project dashboard summary",
	Long:  "Display a summary of all items grouped by status.",
	Args:  cobra.NoArgs,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	s, err := getStore()
	if err != nil {
		return err
	}

	cfg, err := s.ReadConfig()
	if err != nil {
		return err
	}

	items, err := s.ListItems()
	if err != nil {
		return err
	}

	f := getFormatter()
	f.StatusSummary(items, cfg.Project)
	return nil
}
