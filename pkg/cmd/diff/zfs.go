package diff

import (
	"bufio"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/andreasgerstmayr/multidiff/pkg/cmd/cli"
)

func pathOfSnapshot(path string, snapshot string) string {
	if snapshot == "" {
		return path
	}

	spl := strings.Split(snapshot, "@")
	fs := spl[0]
	snap := spl[1]
	return strings.Replace(path, fs, fmt.Sprintf("%s/.zfs/snapshot/%s", fs, snap), 1)
}

func diffSnapshots(opts cli.Options, snapshot1 string, snapshot2 string) ([]Diff, error) {
	args := []string{
		"zfs",
		"diff",
		"-H", // do not print header
		"-h", // do not escape non-ASCII
		"-F", // show file type (file or directory, etc.)
		snapshot1,
	}
	if snapshot2 != "" {
		args = append(args, snapshot2)
	}
	cmd := exec.Command(args[0], args[1:]...)
	log.V(2).Info("running command", "cmd", strings.Join(args, " "))

	out, err := cmd.StdoutPipe()
	if err != nil {
		return []Diff{}, err
	}

	err = cmd.Start()
	if err != nil {
		return []Diff{}, err
	}

	diffMap := map[string]Diff{}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		spl := strings.Split(line, "\t")
		change := spl[0]
		type_ := spl[1]
		source := spl[2]
		dest := ""

		if type_ != "F" {
			continue
		} else if !matchPatterns(opts.IncludePatterns, opts.ExcludePatterns, source) {
			continue
		} else if change == "R" { // renamed
			dest = spl[3]
			if !matchPatterns(opts.IncludePatterns, opts.ExcludePatterns, dest) {
				continue
			}
		}

		lastDiff, exists := diffMap[source]
		if exists {
			log.V(2).Info("found duplicate occurrence of file in snapshot", "prev", lastDiff.Change, "cur", change, "a", lastDiff.PathA, "b", pathOfSnapshot(source, snapshot2))
		}

		switch change {
		case "M": // modified
			diffMap[source] = Diff{
				Change: Modified,
				Source: source,
				PathA:  pathOfSnapshot(source, snapshot1),
				PathB:  pathOfSnapshot(source, snapshot2),
			}

		case "+": // created
			if exists {
				// sometimes applications are updating a file by deleting and creating it
				// merge them into a single, modified event
				diffMap[source] = Diff{
					Change: Modified,
					Source: source,
					PathA:  pathOfSnapshot(source, snapshot1),
					PathB:  pathOfSnapshot(source, snapshot2),
				}
			} else {
				diffMap[source] = Diff{
					Change: Created,
					Source: source,
					PathB:  pathOfSnapshot(source, snapshot2),
				}
			}

		case "-": // removed
			diffMap[source] = Diff{
				Change: Removed,
				Source: source,
				PathA:  pathOfSnapshot(source, snapshot1),
			}

		case "R": // renamed
			diffMap[source] = Diff{
				Change: Renamed,
				Source: source,
				Dest:   dest,
				PathA:  pathOfSnapshot(source, snapshot1),
				PathB:  pathOfSnapshot(dest, snapshot2),
			}

		default:
			continue
		}
	}

	sources := make([]string, 0, len(diffMap))
	for source := range diffMap {
		sources = append(sources, source)
	}
	sort.Strings(sources)

	diffs := []Diff{}
	for _, source := range sources {
		diff := diffMap[source]

		switch diff.Change {
		case Modified, Renamed:
			if !isModified(opts, diff.PathA, diff.PathB) {
				continue
			}
		}

		diffs = append(diffs, diff)
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return diffs, nil
}
