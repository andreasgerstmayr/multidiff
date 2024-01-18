package zgit

import (
	"github.com/spf13/cobra"
)

func NewZgitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zgit",
		Short: "zgit allows exploring ZFS snapshots like git repositories",
	}
	cmd.AddCommand(NewLogCommand())
	cmd.AddCommand(NewShowCommand())
	cmd.AddCommand(NewDiffCommand())
	cmd.AddCommand(NewStatusCommand())
	return cmd
}
