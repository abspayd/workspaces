package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:   "rm path",
		Short: "Remove an existing workspace directory entry",
		Long: `Remove an existing workspace directory entry. This will prevent
				workspaces from opening or listing this workspace.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			for _, arg := range args {
				delete(workspaces.Paths, arg)
				workspaces_dirty = true

				fmt.Printf("Removed \"%s\" from registered workspaces.\n", arg)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(rmCmd)
}
