package zgit

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Show(opts Options, snapshot string) error {
	snapshots, err := listSnapshotsAtCwd()
	if err != nil {
		return err
	}

	found := false
	for i := range snapshots {
		if snapshots[i].Name == snapshot || snapshot == "" {
			printSnapshotHeader(snapshots[i])
			err = Diff(opts, snapshots[i+1].Name, snapshots[i].Name)
			if err != nil {
				return err
			}

			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("snapshot %s not found", snapshot)
	}

	return nil
}

func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [snapshot]",
		Short: "show changes of a snapshot",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmd.Context().Value(Options{}).(Options)
			if len(args) == 1 {
				return Show(opts, args[0])
			} else {
				return Show(opts, "")
			}
		},
	}
	return cmd
}
