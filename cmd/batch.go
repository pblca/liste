package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Execute multiple commands from stdin",
	Long: `Reads commands from stdin (one per line) and executes them sequentially.
Useful for AI agents to emit a batch of mutations in one shot.

Each line is a liste command without the 'liste' prefix.
Lines starting with # are comments and are skipped.
Empty lines are skipped.

Example input:
  add feature "User authentication" --phase 1 --priority high
  add task "Write tests for auth" --phase 1
  link FEAT-001 parent-of TASK-001
  set FEAT-001 status active`,
	Args: cobra.NoArgs,
	RunE: runBatch,
}

func init() {
	rootCmd.AddCommand(batchCmd)
}

func runBatch(cmd *cobra.Command, args []string) error {
	scanner := bufio.NewScanner(os.Stdin)
	lineNum := 0
	executed := 0
	errors := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the line into args (respecting quoted strings)
		cmdArgs, err := splitArgs(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line %d: parse error: %s\n", lineNum, err)
			errors++
			continue
		}

		// Execute as a subcommand of root
		rootCmd.SetArgs(cmdArgs)
		if err := rootCmd.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "line %d: %s\n", lineNum, err)
			errors++
		} else {
			executed++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	if flagJSON {
		fmt.Fprintf(os.Stdout, `{"executed": %d, "errors": %d, "total_lines": %d}`+"\n", executed, errors, lineNum)
	} else if !flagQuiet {
		fmt.Fprintf(os.Stdout, "\nBatch complete: %d executed, %d errors, %d lines read\n", executed, errors, lineNum)
	}

	if errors > 0 {
		return fmt.Errorf("%d command(s) failed", errors)
	}
	return nil
}

// splitArgs splits a command line string into arguments, respecting quoted strings.
func splitArgs(line string) ([]string, error) {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := byte(0)

	for i := 0; i < len(line); i++ {
		ch := line[i]

		if inQuote {
			if ch == quoteChar {
				inQuote = false
			} else {
				current.WriteByte(ch)
			}
		} else {
			switch ch {
			case '"', '\'':
				inQuote = true
				quoteChar = ch
			case ' ', '\t':
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			default:
				current.WriteByte(ch)
			}
		}
	}

	if inQuote {
		return nil, fmt.Errorf("unclosed quote")
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args, nil
}
