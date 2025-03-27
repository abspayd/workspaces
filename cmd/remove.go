package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:   "rm path",
		Short: "Remove an existing workspace directory entry",
		Long: `Remove an existing workspace directory entry. This will prevent
				workspaces from opening or listing this workspace.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Expand any relative path to absolute path
			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			for i, w := range workspace_layout.Workspaces {
				if w.Path == path {
					workspace_layout.Workspaces = append(workspace_layout.Workspaces[:i], workspace_layout.Workspaces[i+1:]...)
					break
				}
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(rmCmd)
}
