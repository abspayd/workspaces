package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Flag variables
	name string

	// Command structure
	addCmd = &cobra.Command{
		Use:   "add [-n workspace name] path",
		Short: "Add a new workspace directory",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			addWorkspace(name, args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the new workspace")
}

func addWorkspace(name string, path string) {
	// TODO
	file, err := os.Open(path)
	if err != nil {
		os.Exit(1)
	}
	_ = file

	logger.Write([]byte(fmt.Sprintf("{[\"%s\", \"%s\"]}")))

	fmt.Printf("Adding workspace with name=%s, path=%s\n", name, path)
}
