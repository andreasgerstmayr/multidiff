package zgit

import (
	"fmt"
	"os"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/andreasgerstmayr/zgit/pkg/logging"
	"github.com/andreasgerstmayr/zgit/pkg/zfs"
)

type Options struct {
	Verbosity     int
	DiffTool      string
	MaxDiffCount  int
	IgnorePattern string
}

const (
	ColorYellow = "\033[93m"
	ColorEnd    = "\033[0m"
)

var log = logging.Logger

func printSnapshotHeader(snapshot zfs.Snapshot) {
	fmt.Printf("%ssnapshot %s%s\n", ColorYellow, snapshot.Name, ColorEnd)
	fmt.Printf("Date:    %s\n", snapshot.Creation.Format("Mon Jan 02 15:04:05 2006 -0700"))
	fmt.Println()
}

func listSnapshotsAtCwd() ([]zfs.Snapshot, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return []zfs.Snapshot{}, err
	}

	fs, err := zfs.FindFilesystem(pwd)
	if err != nil {
		return []zfs.Snapshot{}, err
	}

	snapshots, err := zfs.ListSnapshots(fs)
	if err != nil {
		return []zfs.Snapshot{}, err
	}

	return snapshots, nil
}

// sometimes zfs diff returns a removal and creation of the same file in the same snapshot
// mergeDiffs merges these into a single diff with type=modified
func mergeDiffs(diffs []zfs.SnapshotDiff) []zfs.SnapshotDiff {
	files := map[string]int{}

	i := 0
	for _, diff := range diffs {
		skip := false

		if diff.Type == "F" {
			lastDiffIdx, ok := files[diff.Source]
			if !ok {
				files[diff.Source] = i
			} else {
				lastDiff := diffs[lastDiffIdx]
				log.V(2).Info("found duplicate occurrence of file in snapshot", "file", diff.Source, "prev", lastDiff.Change, "cur", diff.Change)
				if diff.Change == "+" && lastDiff.Change == "-" {
					lastDiff.Change = "M"
					diffs[lastDiffIdx] = lastDiff
					skip = true
				}
			}
		}

		if !skip {
			diffs[i] = diff
			i++
		}
	}
	diffs = diffs[:i]
	return diffs
}

func filterDiffs(diffs []zfs.SnapshotDiff, ignorePattern string) []zfs.SnapshotDiff {
	if ignorePattern == "" {
		return diffs
	}

	i := 0
	for _, diff := range diffs {
		skip := false

		if diff.Type == "F" {
			match, err := doublestar.Match(ignorePattern, diff.Source)
			if err != nil {
				log.Error(err, "error matching file pattern")
			} else {
				skip = match
			}
		}

		if !skip {
			diffs[i] = diff
			i++
		}
	}
	diffs = diffs[:i]
	return diffs
}
