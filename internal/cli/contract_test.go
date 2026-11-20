package cli

import (
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// TestEveryCommandHonoursTheCLIContract walks the whole command tree, so every new
// command is held to the same documentation standard without anyone adding a test.
func TestEveryCommandHonoursTheCLIContract(t *testing.T) {
	for _, cmd := range allCommands(NewRootCommand(io.Discard, io.Discard)) {
		t.Run(cmd.CommandPath(), func(t *testing.T) {
			assertDocumented(t, cmd)
			assertFlagsDescribed(t, cmd)
		})
	}
}

func assertDocumented(t *testing.T, cmd *cobra.Command) {
	t.Helper()
	if strings.TrimSpace(cmd.Short) == "" {
		t.Error("missing Short description")
	}
	if cmd.Name() != strings.ToLower(cmd.Name()) {
		t.Errorf("command name %q must be lowercase", cmd.Name())
	}
	if cmd.Runnable() && strings.TrimSpace(cmd.Example) == "" {
		t.Error("runnable command is missing an Example")
	}
}

func assertFlagsDescribed(t *testing.T, cmd *cobra.Command) {
	t.Helper()
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if strings.TrimSpace(flag.Usage) == "" {
			t.Errorf("flag --%s has no usage text", flag.Name)
		}
	})
}

func allCommands(cmd *cobra.Command) []*cobra.Command {
	commands := []*cobra.Command{cmd}
	for _, child := range cmd.Commands() {
		commands = append(commands, allCommands(child)...)
	}
	return commands
}
