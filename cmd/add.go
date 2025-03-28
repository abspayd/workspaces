package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	stow bool

	// Command structure
	addCmd = &cobra.Command{
		Use:   "add path",
		Short: "Add a new workspace directory",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				path, err := filepath.Abs(arg)
				if err != nil {
					log.Println("Failed to resolve path:", err)
					return err
				}
				ws := Workspace{
					StowDir: stow,
				}

				workspaces.Paths[path] = ws
				workspaces_dirty = true

				fmt.Printf("Added \"%s\" to registered workspaces.\n", path)
			}

			return nil
		},
	}
)

func init() {
	addCmd.Flags().BoolVarP(&stow, "stow", "s", false, "Add directories as GNU Stow directories. This will allow project paths to follow the stow package structure.")
	rootCmd.AddCommand(addCmd)
}
