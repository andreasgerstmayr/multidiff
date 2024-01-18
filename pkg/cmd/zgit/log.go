package zgit

import (
	"fmt"

	"github.com/andreasgerstmayr/multidiff/pkg/cmd/cli"
	"github.com/andreasgerstmayr/multidiff/pkg/cmd/diff"
	"github.com/spf13/cobra"
)

type LogOptions struct {
	cli.Options
	ShowPatch  bool
	IgnoreMeta bool
}

func Log(opts LogOptions) error {
	snapshots, err := listSnapshotsAtCwd()
	if err != nil {
		return err
	}

	for i := 0; i < len(snapshots)-1; i++ {
		if snapshots[i].Used == "0B" {
			continue
		}

		prevSnapshot := snapshots[i+1].Name
		curSnapshot := snapshots[i].Name

		diffs, err := diff.List(opts.Options, prevSnapshot, curSnapshot)
		if err != nil {
			return err
		}

		if len(diffs) == 0 {
			continue
		}

		printSnapshotHeader(snapshots[i])
		if opts.ShowPatch {
			diff.Print(opts.Options, diffs)
			fmt.Println()
		}
	}

	return nil
}

func NewLogCommand() *cobra.Command {
	opts := LogOptions{}
	cmd := &cobra.Command{
		Use:   "log",
		Short: "list all non empty snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Options = cmd.Context().Value(cli.OptionsKey{}).(cli.Options)
			return Log(opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.ShowPatch, "patch", "p", false, "show patch")
	return cmd
}
