package zfs

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/andreasgerstmayr/zgit/pkg/logging"
)

type Snapshot struct {
	Name     string
	Used     string
	Creation time.Time
}

type SnapshotDiff struct {
	Change string
	Type   string
	Source string
	Dest   string
}

var log = logging.Logger

func FindFilesystem(path string) (string, error) {
	args := []string{path}
	cmd := exec.Command("df", args...)
	log.V(2).Info("running command", "cmd", "df "+strings.Join(args, " "))

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	s := string(out)
	spl := strings.Split(s, "\n")
	if len(spl) != 3 {
		return "", fmt.Errorf("invalid output from df: %s", s)
	}

	fs := strings.Split(spl[1], " ")[0]
	return fs, nil
}

func GetFileOfSnapshot(path string, snapshot string) string {
	spl := strings.Split(snapshot, "@")
	fs := spl[0]
	snap := spl[1]
	return strings.Replace(path, fs, fmt.Sprintf("%s/.zfs/snapshot/%s", fs, snap), 1)
}

func ListSnapshots(fs string) ([]Snapshot, error) {
	args := []string{
		"list",
		"-H",             // do not print header
		"-t", "snapshot", // show snapshots
		"-o", "name,used,creation",
		"-S", "creation", // sort by creation
		fs,
	}
	cmd := exec.Command("zfs", args...)
	log.V(2).Info("running command", "cmd", "zfs "+strings.Join(args, " "))

	out, err := cmd.StdoutPipe()
	if err != nil {
		return []Snapshot{}, err
	}

	err = cmd.Start()
	if err != nil {
		return []Snapshot{}, err
	}

	snapshots := []Snapshot{}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		spl := strings.Split(line, "\t")
		creation, err := time.ParseInLocation("Mon Jan 2  15:04 2006", spl[2], time.Local)
		if err != nil {
			return []Snapshot{}, err
		}

		snapshots = append(snapshots, Snapshot{
			Name:     spl[0],
			Used:     spl[1],
			Creation: creation,
		})
	}
	return snapshots, nil
}

func DiffSnapshots(snapshot1 string, snapshot2 string, maxDiffCount int) ([]SnapshotDiff, error) {
	args := []string{
		"diff",
		"-H", // do not print header
		"-h", // do not escape non-ASCII
		"-F", // show file type
		snapshot1,
	}
	if snapshot2 != "" {
		args = append(args, snapshot2)
	}
	cmd := exec.Command("zfs", args...)
	log.V(2).Info("running command", "cmd", "zfs "+strings.Join(args, " "))

	out, err := cmd.StdoutPipe()
	if err != nil {
		return []SnapshotDiff{}, err
	}

	err = cmd.Start()
	if err != nil {
		return []SnapshotDiff{}, err
	}

	diffs := []SnapshotDiff{}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		spl := strings.Split(line, "\t")
		switch spl[0] {
		case "R": // renamed
			diffs = append(diffs, SnapshotDiff{
				Change: spl[0],
				Type:   spl[1],
				Source: spl[2],
				Dest:   spl[3],
			})
		case "M", "+", "-": // modified, created, removed
			diffs = append(diffs, SnapshotDiff{
				Change: spl[0],
				Type:   spl[1],
				Source: spl[2],
			})
		}

		if maxDiffCount != -1 && len(diffs) > maxDiffCount {
			log.Info("skipping diffs: reached maximum number of diffs per snapshot", "snapshot1", snapshot1, "snapshot2", snapshot2, "maxDiffCount", maxDiffCount)
			// return empty set because returning only the first X diffs could result in incorrect diffs,
			// if the same file appears multiple times in the same diff (see mergeDiffs)
			return []SnapshotDiff{}, nil
		}
	}
	return diffs, nil
}
