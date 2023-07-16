package main

import (
	"context"
	"os"

	"github.com/andreasgerstmayr/zgit/pkg/cmd/zgit"
	"github.com/andreasgerstmayr/zgit/pkg/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

func newRootCommand() *cobra.Command {
	opts := zgit.Options{}
	cmd := &cobra.Command{
		Use:          "zgit",
		Long:         "zgit allows exploring ZFS snapshots like git repositories",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logging.Level.SetLevel(zapcore.Level(-opts.Verbosity))
			cmd.SetContext(context.WithValue(cmd.Context(), zgit.Options{}, opts))
		},
	}
	cmd.PersistentFlags().CountVarP(&opts.Verbosity, "verbose", "v", "verbosity")
	cmd.PersistentFlags().StringVar(&opts.DiffTool, "difftool", "git --no-pager diff --no-index --color=always", "use a custom diff program")
	cmd.PersistentFlags().IntVar(&opts.MaxDiffCount, "max-diff-count", 50, "maximum number of changes per snapshot")
	cmd.PersistentFlags().StringVarP(&opts.IgnorePattern, "ignore", "i", "", "ignore files matching this pattern in the diff")
	return cmd
}

func main() {
	rootCmd := newRootCommand()
	rootCmd.AddCommand(zgit.NewLogCommand())
	rootCmd.AddCommand(zgit.NewShowCommand())
	rootCmd.AddCommand(zgit.NewDiffCommand())
	rootCmd.AddCommand(zgit.NewStatusCommand())

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
