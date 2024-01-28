package diff

import (
	"os"

	"github.com/andreasgerstmayr/multidiff/pkg/cmd/cli"
	"github.com/andreasgerstmayr/multidiff/pkg/logging"
)

var log = logging.Logger

func List(opts cli.Options, a string, b string) ([]Diff, error) {
	switch {
	case isFile(a) && isFile(b):
		return diffFiles(opts, a, b)
	default:
		return diffSnapshots(opts, a, b)
	}
}

func diffFiles(opts cli.Options, a string, b string) ([]Diff, error) {
	if isModified(opts, a, b) {
		return []Diff{{
			Change: Modified,
			Source: a,
			PathA:  a,
			PathB:  b,
		}}, nil
	} else {
		return []Diff{}, nil
	}
}
func isFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isModified(opts cli.Options, a string, b string) bool {
	eq, err := compareFileSize(a, b)
	if err != nil {
		log.Error(err, "error checking file size", "a", a, "b", b)
	} else if !eq {
		// if the file size is not equal, file is modified
		return true
	}

	if !opts.ShowMetadataChanges {
		eq, err := compareModificationTime(a, b)
		if err != nil {
			log.Error(err, "error checking modified time", "a", a, "b", b)
		} else if eq {
			// if modification time is equal, file content was not changed
			return false
		}
	}

	if opts.CompareByteForByte {
		eq, err := compareFilesByteForByte(a, b, 40*1024)
		if err != nil {
			log.Error(err, "error comparing files byte-for-byte", "a", a, "b", b)
		} else if eq {
			// file content was not changed
			return false
		}
	}

	return true
}
