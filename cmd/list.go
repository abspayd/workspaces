package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "ls",
		Short: "List all workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cobra.NoArgs(cmd, args)
			if err != nil {
				return err
			}

			for workspace := range workspaces.Paths {
				if workspaces.Paths[workspace].StowDir {
					fmt.Println(workspace, "-", "stow")
				} else {
					fmt.Println(workspace)
				}
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}
