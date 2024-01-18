package diff

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/andreasgerstmayr/multidiff/pkg/cmd/cli"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

func Show(opts cli.Options, a string, b string) error {
	diffs, err := List(opts, a, b)
	if err != nil {
		log.Error(err, "error getting diffs", "a", a, "b", b)
		os.Exit(1)
	}

	Print(opts, diffs)
	return nil
}

func Print(opts cli.Options, diffs []Diff) {
	for _, d := range diffs {
		if opts.ShowPathOnly {
			PrintPath(d)
		} else {
			err := PrintDiff(opts, d)
			if err != nil {
				log.Error(err, "error printing diff", "a", d.PathA, "b", d.PathB)
			}
		}
	}
}

func PrintPath(diff Diff) {
	switch diff.Change {
	case Modified:
		fmt.Printf("M\t%s\n", diff.Source)

	case Created:
		fmt.Printf("A\t%s\n", diff.Source)

	case Removed:
		fmt.Printf("D\t%s\n", diff.Source)

	case Renamed:
		fmt.Printf("R\t%s\t%s\n", diff.Source, diff.Dest)
	}
}

func PrintDiff(opts cli.Options, diff Diff) error {
	bold := color.New(color.Bold)

	switch diff.Change {
	case Modified:
		bold.Printf("diff --git a%s b%s\n", diff.PathA, diff.PathB)
		return convDiff(opts, diff.PathA, diff.PathB)

	case Renamed:
		bold.Printf("diff --git a%s b%s\n", diff.Source, diff.Dest)
		bold.Printf("rename from %s\n", diff.Source)
		bold.Printf("rename to   %s\n", diff.Dest)
		return convDiff(opts, diff.PathA, diff.PathB)

	case Created:
		if opts.NewFilesAsEmpty {
			bold.Printf("diff --git a%s b%s\n", "/dev/null", diff.PathB)
			return gitDiff("/dev/null", diff.PathB)
		} else {
			bold.Printf("Created: %s\n", diff.Source)
			return nil
		}

	case Removed:
		if opts.NewFilesAsEmpty {
			bold.Printf("diff --git a%s b%s\n", diff.PathA, "/dev/null")
			return gitDiff(diff.PathA, "/dev/null")
		} else {
			bold.Printf("Removed: %s\n", diff.Source)
			return nil
		}

	default:
		return nil
	}
}

func gitDiff(a string, b string) error {
	args := []string{
		"git",
		"diff",
		"--no-index", // do not require a git repository
		a,
		b,
	}
	if isatty.IsTerminal(os.Stdout.Fd()) {
		args = append(args, "--color=always")
	}
	cmd := exec.Command(args[0], args[1:]...)
	log.V(2).Info("running command", "cmd", strings.Join(args, " "))

	cmd.Stderr = os.Stderr
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(out)
	scanner.Scan() // skip 'diff --git' line
	scanner.Scan() // skip 'index' line
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	cmd.Wait() // git diff --no-index exits with 1 if there were changes
	return nil
}

func NewDiffCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "diff a [b]",
		Short: "show diff between two sources",
		Long:  "show diff between two sources (if only given one argument, show diff between source and working copy)",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmd.Context().Value(cli.OptionsKey{}).(cli.Options)
			if len(args) == 1 {
				return Show(opts, args[0], "")
			} else {
				return Show(opts, args[0], args[1])
			}
		},
	}
}
