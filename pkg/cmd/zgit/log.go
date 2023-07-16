package zgit

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/andreasgerstmayr/zgit/pkg/zfs"
	"github.com/spf13/cobra"
)

type LogOptions struct {
	Global     Options
	ShowPatch  bool
	IgnoreMeta bool
}

// compareFiles checks if file1 and file2 have the same content
func compareFiles(file1 string, file2 string, chunkSize int) (bool, error) {
	// first compare file modification time
	st1, err := os.Stat(file1)
	if err != nil {
		return false, err
	}

	st2, err := os.Stat(file2)
	if err != nil {
		return false, err
	}

	if st1.ModTime().Equal(st2.ModTime()) {
		return true, nil
	}

	// then compare file contents
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		n1, err1 := f1.Read(b1)
		if err1 != nil && err1 != io.EOF {
			return false, err1
		}

		b2 := make([]byte, chunkSize)
		n2, err2 := f2.Read(b2)
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}

		// exit if we're at EOF of a file
		if err1 == io.EOF && err2 == io.EOF {
			return true, nil
		} else if err1 == io.EOF || err2 == io.EOF {
			return false, nil
		}

		if n1 != n2 || !bytes.Equal(b1, b2) {
			return false, nil
		}
	}
}

func containsFileContentChange(diffs []zfs.SnapshotDiff, snapshot1 string, snapshot2 string) (bool, error) {
	for _, diff := range diffs {
		if diff.Type != "F" {
			continue
		}

		if diff.Change == "+" || diff.Change == "-" {
			return true, nil
		} else if diff.Change == "M" {
			file1 := zfs.GetFileOfSnapshot(diff.Source, snapshot1)
			var file2 string
			if snapshot2 != "" {
				file2 = zfs.GetFileOfSnapshot(diff.Source, snapshot2)
			} else {
				file2 = diff.Source
			}

			cmp, err := compareFiles(file1, file2, 40*1024)
			if err != nil {
				return false, err
			}

			log.V(2).Info("compared files", "file1", file1, "file2", file2, "equal", cmp)
			if !cmp {
				return true, nil
			}
		}
	}
	return false, nil
}

func Log(opts LogOptions) error {
	snapshots, err := listSnapshotsAtCwd()
	if err != nil {
		return err
	}

	for i := range snapshots {
		if snapshots[i].Used == "0B" {
			continue
		}

		prevSnapshot := snapshots[i+1].Name
		curSnapshot := snapshots[i].Name

		diffs, err := zfs.DiffSnapshots(prevSnapshot, curSnapshot, opts.Global.MaxDiffCount)
		if err != nil {
			return err
		}

		diffs = mergeDiffs(diffs)
		diffs = filterDiffs(diffs, opts.Global.IgnorePattern)

		if len(diffs) == 0 {
			continue
		}

		if opts.IgnoreMeta {
			fileContentChange, err := containsFileContentChange(diffs, prevSnapshot, curSnapshot)
			if err != nil {
				return err
			}

			if !fileContentChange {
				log.V(1).Info("skipping snapshot because it contains no file content change", "snapshot", curSnapshot)
				continue
			}
		}

		printSnapshotHeader(snapshots[i])
		if opts.ShowPatch {
			err = printDiff(opts.Global, diffs, prevSnapshot, curSnapshot)
			if err != nil {
				return err
			}

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
			opts.Global = cmd.Context().Value(Options{}).(Options)
			return Log(opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.ShowPatch, "patch", "p", false, "show patch")
	cmd.Flags().BoolVarP(&opts.IgnoreMeta, "ignore-meta", "m", true, "ignore metadata modifications")
	return cmd
}
