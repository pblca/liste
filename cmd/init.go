package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pblca/liste/internal/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new .liste/ in the current directory",
	Long:  "Creates a .liste/ directory with default config and state files.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	name := filepath.Base(cwd)
	if len(args) > 0 {
		name = args[0]
	}

	roadmapPath := filepath.Join(cwd, ".liste")
	s := store.New(roadmapPath)

	if s.Exists() {
		return fmt.Errorf(".liste/ already exists in %s", cwd)
	}

	if err := s.Init(name); err != nil {
		return err
	}

	f := getFormatter()
	f.Message(fmt.Sprintf("Initialized .liste/ for project %q", name))
	return nil
}
