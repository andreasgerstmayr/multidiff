package zgit

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/andreasgerstmayr/multidiff/pkg/logging"
)

const (
	ColorYellow = "\033[93m"
	ColorEnd    = "\033[0m"
)

type Snapshot struct {
	Name     string
	Used     string
	Creation time.Time
}

var log = logging.Logger

func printSnapshotHeader(snapshot Snapshot) {
	fmt.Printf("%ssnapshot %s%s\n", ColorYellow, snapshot.Name, ColorEnd)
	fmt.Printf("Date:    %s\n", snapshot.Creation.Format("Mon Jan 02 15:04:05 2006 -0700"))
	fmt.Println()
}

func listSnapshotsAtCwd() ([]Snapshot, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return []Snapshot{}, err
	}

	fs, err := findFilesystem(pwd)
	if err != nil {
		return []Snapshot{}, err
	}

	snapshots, err := listSnapshots(fs)
	if err != nil {
		return []Snapshot{}, err
	}

	return snapshots, nil
}

func findFilesystem(path string) (string, error) {
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

func listSnapshots(fs string) ([]Snapshot, error) {
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
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	snapshots := []Snapshot{}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		spl := strings.Split(line, "\t")
		creation, err := time.ParseInLocation("Mon Jan 2  15:04 2006", spl[2], time.Local)
		if err != nil {
			return nil, err
		}

		snapshots = append(snapshots, Snapshot{
			Name:     spl[0],
			Used:     spl[1],
			Creation: creation,
		})
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}
