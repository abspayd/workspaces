package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// Command structure
	addCmd = &cobra.Command{
		Use:   "add path",
		Short: "Add a new workspace directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Expand any relative path to absolute path
			path, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}

			for _, workspace := range workspace_layout.Workspaces {
				if workspace.Path == path {
					cmd.Printf("Workspace \"%s\" already exists.\n", path)
					return nil
				}
			}

			workspace := Workspace{
				Path:    path,
				StowDir: false,
			}
			workspace_layout.Workspaces = append(workspace_layout.Workspaces, workspace)

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
