package diff

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/andreasgerstmayr/multidiff/pkg/cmd/cli"
)

// optionally convert files to a textual representation,
// then runs git diff
func convDiff(opts cli.Options, a string, b string) error {
	extA := strings.ToLower(path.Ext(a))
	extB := strings.ToLower(path.Ext(b))

	if extA == extB {
		switch extA {
		case ".jpg", ".heic":
			err := exiftool(a, "/tmp/multidiff_a.txt")
			if err != nil {
				return err
			}
			a = "/tmp/multidiff_a.txt"

			err = exiftool(b, "/tmp/multidiff_b.txt")
			if err != nil {
				return err
			}
			b = "/tmp/multidiff_b.txt"
		}
	}

	return gitDiff(a, b)
}

func exiftool(src string, dest string) error {
	args := []string{
		"exiftool",
		"-a",         // show duplicate tag names
		"-u",         // show unknown tags
		"-G1",        // sort by group
		"--File:all", // do not print file information (filename, permissions, etc)
		src,
	}
	cmd := exec.Command(args[0], args[1:]...)
	log.V(2).Info("running command", "cmd", strings.Join(args, " "))

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd.Stdout = file
	return cmd.Run()
}
