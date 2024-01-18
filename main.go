package main

import (
	"context"
	"os"

	"github.com/andreasgerstmayr/multidiff/pkg/cmd/cli"
	"github.com/andreasgerstmayr/multidiff/pkg/cmd/diff"
	"github.com/andreasgerstmayr/multidiff/pkg/cmd/zgit"
	"github.com/andreasgerstmayr/multidiff/pkg/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

func newRootCommand() *cobra.Command {
	opts := cli.Options{}
	cmd := &cobra.Command{
		Use:          "multidiff",
		Long:         "multidiff compares two sources (files/directories/ZFS snapshots) in a human readable way (e.g. metadata of images)",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logging.Level.SetLevel(zapcore.Level(-opts.Verbosity))
			cmd.SetContext(context.WithValue(cmd.Context(), cli.OptionsKey{}, opts))
		},
	}
	cmd.PersistentFlags().CountVarP(&opts.Verbosity, "verbose", "v", "set verbosity (can be used multiple times)")
	cmd.PersistentFlags().BoolVar(&opts.ShowPathOnly, "path", false, "show paths only")
	cmd.PersistentFlags().BoolVarP(&opts.NewFilesAsEmpty, "new-file", "N", false, "show absent files as empty")
	cmd.PersistentFlags().BoolVarP(&opts.ShowMetadataChanges, "show-meta", "m", false, "include file metadata modifications in diff")
	cmd.PersistentFlags().BoolVarP(&opts.CompareByteForByte, "byte", "b", true, "compare byte for byte")
	cmd.PersistentFlags().StringArrayVarP(&opts.IncludePatterns, "include", "i", []string{}, "include files matching this pattern")
	cmd.PersistentFlags().StringArrayVarP(&opts.ExcludePatterns, "exclude", "e", []string{}, "exclude files matching this pattern")
	cmd.PersistentFlags().BoolVar(&opts.Conv.Exif, "conv.exif", false, "compare EXIF metadata of images")

	return cmd
}

func main() {
	rootCmd := newRootCommand()
	rootCmd.AddCommand(diff.NewDiffCommand())
	rootCmd.AddCommand(zgit.NewZgitCommand())

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
