package zgit

import (
	"github.com/andreasgerstmayr/multidiff/pkg/cmd/cli"
	"github.com/andreasgerstmayr/multidiff/pkg/cmd/diff"
	"github.com/spf13/cobra"
)

func NewDiffCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "diff [snapshot] [snapshot2]",
		Short: "show diff between two snapshots",
		Long:  "show diff between two snapshots (if only given one argument, show diff between snapshot and working copy)",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmd.Context().Value(cli.OptionsKey{}).(cli.Options)
			if len(args) == 1 {
				return diff.Show(opts, args[0], "")
			} else {
				return diff.Show(opts, args[0], args[1])
			}
		},
	}
}
