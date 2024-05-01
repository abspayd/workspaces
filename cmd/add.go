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

			for i := 0; i < len(workspaces); i++ {
				if workspaces[i] == path {
					// This path already exists; do nothing
					return nil
				}
			}

			workspaces = append(workspaces, path)

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}
