package zgit

import (
	"github.com/spf13/cobra"
)

func Status(opts Options) error {
	snapshots, err := listSnapshotsAtCwd()
	if err != nil {
		return err
	}

	err = Diff(opts, snapshots[0].Name, "")
	if err != nil {
		return err
	}

	return nil
}

func NewStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "show changes in working copy",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmd.Context().Value(Options{}).(Options)
			return Status(opts)
		},
	}
	return cmd
}
