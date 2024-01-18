package zgit

import (
	"github.com/andreasgerstmayr/multidiff/pkg/cmd/cli"
	"github.com/andreasgerstmayr/multidiff/pkg/cmd/diff"
	"github.com/spf13/cobra"
)

func Status(opts cli.Options) error {
	snapshots, err := listSnapshotsAtCwd()
	if err != nil {
		return err
	}

	err = diff.Show(opts, snapshots[0].Name, "")
	if err != nil {
		return err
	}

	return nil
}

func NewStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "show changes in working copy",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmd.Context().Value(cli.OptionsKey{}).(cli.Options)
			return Status(opts)
		},
	}
}
