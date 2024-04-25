package cmd

import (
	"os"

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
	logger  *os.File
	rootCmd = &cobra.Command{
		Use:   "workspaces",
		Short: "Workspaces is a simple CLI workspace manager.",
		Long:  `A simple workspace manager to quickly access and modify your frequent project locations.`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// cobra.OnInitialize() @TODO: do I want a config file? Probably not, for now...
	home, err := os.UserHomeDir()
	if err != nil {
		rootCmd.PrintErrln("Error: Unable to locate home user directory")
		os.Exit(1)
	}
	logger_path := home + "/.local/share/workspaces"
	workspace_log := logger_path + "workspaces.json"
	logger, err := os.OpenFile(workspace_log, os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		rootCmd.PrintErrln("Error: Unable to access workspaces.")
		os.Exit(1)
	}

	var buf []byte
	n, err := logger.Read(buf)
	if n == 0 || err != nil {
	}

	cobra.OnFinalize(cleanup)
}

func cleanup() {
	logger.Close()
}
