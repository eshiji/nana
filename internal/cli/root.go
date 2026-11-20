// Package cli maps command-line commands onto the application. It holds no business logic.
package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// Exit codes are part of the CLI contract: 0 success, 1 findings, 2 usage or runtime error.
const (
	ExitOK    = 0
	ExitError = 2
)

// Execute runs the CLI with the given arguments and returns the process exit code.
func Execute(args []string, stdout, stderr io.Writer) int {
	root := NewRootCommand(stdout, stderr)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(stderr, "error:", err)
		return ExitError
	}
	return ExitOK
}

// NewRootCommand builds the full command tree. Execute and tools/gendocs share it.
func NewRootCommand(stdout, stderr io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:           "nana",
		Short:         "Pipeline checks for container images and Terraform, the same on a laptop and in CI.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.AddCommand(newVersionCommand(), newDockerfileCommand())
	return root
}
