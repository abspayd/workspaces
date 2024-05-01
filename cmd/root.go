package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

/*
Command examples:
	workspaces add path
	workspaces remove path
	workspaces list [--filter string]
	workspaces open name
*/

var (
	workspaces []string

	workspaces_list_file string

	rootCmd = &cobra.Command{
		Use:   "workspaces",
		Short: "A simple CLI workspace manager.",
		Long:  `A simple workspace manager to quickly access and modify your frequent project locations.`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: Unable to locate home user directory")
		os.Exit(1)
	}
	workspaces_path := home + "/.local/share/workspaces"
	workspaces_list_file = workspaces_path + "/workspaces.csv"

	// Check if workspace directory exists
	_, err = os.Stat(workspaces_path)
	if err != nil && os.IsNotExist(err) {
		// The directory does not exit. Create the workspace directory path
		if err = os.MkdirAll(workspaces_path, 0700); err != nil {
			fmt.Println("Failed to create workspace list directory.")
			os.Exit(1)
		}
	}

	// Check if workspace list file exists
	buf, err := os.ReadFile(workspaces_list_file)
	if err != nil && os.IsNotExist(err){
		// The file does not exist. Create it and close it.
		fd, err := os.Create(workspaces_list_file)
		if err != nil {
			fmt.Println("Failed to create workspace list directory.")
			os.Exit(1)
		}
		fd.Close()
	}

	if len(buf) > 0 {
		workspaces = strings.Split(string(buf), ",")
	} else {
		workspaces = nil
	}

	cobra.OnFinalize(finalize)
}

func finalize() {
	err := os.WriteFile(workspaces_list_file, []byte(strings.Join(workspaces, ",")), 0700)
	if err != nil {
		fmt.Println("Error: Failed to add changes.")
		os.Exit(1)
	}
}
