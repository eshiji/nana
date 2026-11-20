package cli

import (
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// commandRule is one requirement every command must meet. Adding a requirement means
// adding one entry to commandRules.
type commandRule struct {
	name     string
	violated func(cmd *cobra.Command) bool
}

var commandRules = []commandRule{
	{"have a Short description", func(c *cobra.Command) bool { return blank(c.Short) }},
	{"have a lowercase name", func(c *cobra.Command) bool { return c.Name() != strings.ToLower(c.Name()) }},
	{"have an Example when runnable", func(c *cobra.Command) bool { return c.Runnable() && blank(c.Example) }},
	{"describe every flag", hasUndescribedFlag},
}

// TestEveryCommandHonoursTheCLIContract walks the whole command tree, so every new
// command is held to the same standard without anyone adding a test.
func TestEveryCommandHonoursTheCLIContract(t *testing.T) {
	for _, cmd := range allCommands(NewRootCommand(io.Discard, io.Discard)) {
		for _, rule := range commandRules {
			if rule.violated(cmd) {
				t.Errorf("%s: must %s", cmd.CommandPath(), rule.name)
			}
		}
	}
}

func blank(s string) bool {
	return strings.TrimSpace(s) == ""
}

func hasUndescribedFlag(cmd *cobra.Command) bool {
	undescribed := false
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		undescribed = undescribed || blank(flag.Usage)
	})
	return undescribed
}

func allCommands(cmd *cobra.Command) []*cobra.Command {
	commands := []*cobra.Command{cmd}
	for _, child := range cmd.Commands() {
		commands = append(commands, allCommands(child)...)
	}
	return commands
}
