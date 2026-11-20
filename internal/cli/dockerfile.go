package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/eshiji/nana/internal/dockerfile"
)

const (
	outputTable = "table"
	outputJSON  = "json"
	noDefault   = "-"
)

func newDockerfileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dockerfile",
		Short: "Read facts from a Dockerfile",
	}
	cmd.AddCommand(newDockerfileArgsCommand())
	return cmd
}

func newDockerfileArgsCommand() *cobra.Command {
	var file, output string
	cmd := &cobra.Command{
		Use:   "args",
		Short: "List the ARG declarations of a Dockerfile with their defaults",
		Example: "  nana dockerfile args\n" +
			"  nana dockerfile args --file build/Dockerfile -o json",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			render, err := argsRenderer(output)
			if err != nil {
				return err
			}
			args, err := readArgs(file)
			if err != nil {
				return err
			}
			return render(cmd.OutOrStdout(), args)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "Dockerfile", "path to the Dockerfile")
	cmd.Flags().StringVarP(&output, "output", "o", outputTable, "output format: table or json")
	return cmd
}

type argsRender func(io.Writer, []dockerfile.Arg) error

func argsRenderer(output string) (argsRender, error) {
	switch output {
	case outputTable:
		return renderArgsTable, nil
	case outputJSON:
		return renderArgsJSON, nil
	}
	return nil, fmt.Errorf("unsupported output format %q (want %s or %s)", output, outputTable, outputJSON)
}

func readArgs(path string) ([]dockerfile.Arg, error) {
	f, err := os.Open(path) // #nosec G304 -- reading the file the user names is the command's purpose
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	args, err := dockerfile.ParseArgs(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return args, nil
}

func renderArgsTable(w io.Writer, args []dockerfile.Arg) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "NAME\tDEFAULT\tSTAGE\tLINE")
	for _, arg := range args {
		value := arg.Default
		if !arg.HasDefault {
			value = noDefault
		}
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", arg.Name, value, arg.Stage, arg.Line)
	}
	return tw.Flush()
}

func renderArgsJSON(w io.Writer, args []dockerfile.Arg) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(args)
}
