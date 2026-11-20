package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/eshiji/nana/internal/buildinfo"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the version this binary was built from",
		Example: "  nana version",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), buildinfo.Version())
			return err
		},
	}
}
