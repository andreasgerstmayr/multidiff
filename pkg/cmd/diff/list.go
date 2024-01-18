package diff

import (
	"github.com/andreasgerstmayr/multidiff/pkg/cmd/cli"
	"github.com/andreasgerstmayr/multidiff/pkg/logging"
)

var log = logging.Logger

func List(opts cli.Options, a string, b string) ([]Diff, error) {
	// TODO: check if source is file, directory or snapshot
	return diffSnapshots(opts, a, b)
}
