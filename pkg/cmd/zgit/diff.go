package zgit

import (
	"os"
	"os/exec"
	"strings"

	"github.com/andreasgerstmayr/zgit/pkg/zfs"
	"github.com/spf13/cobra"
)

func runDiffTool(diffTool string, file1 string, file2 string) {
	spl := strings.SplitN(diffTool, " ", 2)
	name := spl[0]
	args := strings.Split(spl[1], " ")
	args = append(args, file1, file2)

	cmd := exec.Command(name, args...)
	log.V(2).Info("running command", "cmd", name+" "+strings.Join(args, " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func printDiff(opts Options, diffs []zfs.SnapshotDiff, snapshot1 string, snapshot2 string) error {
	for _, diff := range diffs {
		log.V(1).Info("diff", "diff", diff)

		if diff.Type == "F" {
			file1 := "/dev/null"
			file2 := "/dev/null"
			switch diff.Change {
			case "M": // file modified
				file1 = zfs.GetFileOfSnapshot(diff.Source, snapshot1)
				if snapshot2 != "" {
					file2 = zfs.GetFileOfSnapshot(diff.Source, snapshot2)
				} else {
					file2 = diff.Source
				}
			case "+": // file created
				if snapshot2 != "" {
					file2 = zfs.GetFileOfSnapshot(diff.Source, snapshot2)
				} else {
					file2 = diff.Source
				}
			case "-": // file removed
				file1 = zfs.GetFileOfSnapshot(diff.Source, snapshot1)
			}

			runDiffTool(opts.DiffTool, file1, file2)
		}
	}
	return nil
}

func Diff(opts Options, snapshot1 string, snapshot2 string) error {
	diffs, err := zfs.DiffSnapshots(snapshot1, snapshot2, opts.MaxDiffCount)
	if err != nil {
		return err
	}

	diffs = mergeDiffs(diffs)
	diffs = filterDiffs(diffs, opts.IgnorePattern)
	return printDiff(opts, diffs, snapshot1, snapshot2)
}

func NewDiffCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff [snapshot] [snapshot2]",
		Short: "show diff between two snapshots (if only given one argument, show diff between snapshot and working copy)",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmd.Context().Value(Options{}).(Options)
			if len(args) == 1 {
				return Diff(opts, args[0], "")
			} else {
				return Diff(opts, args[0], args[1])
			}
		},
	}
	return cmd
}
