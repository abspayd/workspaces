package cmd

import (
	"errors"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	closeCmd = &cobra.Command{
		Use:   "close [name]",
		Short: "Close a project session",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := workspaceProjects()
			if err != nil {
				return err
			}

			// FZF
			fzf_query := ""
			if len(args) > 0 {
				fzf_query = args[0]
			}
			project_name, err := fzf(fzf_query, projects)
			if err != nil {
				exit_error := &exec.ExitError{}
				if errors.As(err, &exit_error) {
					return nil
				}
				return err
			}

			// == Tmux ==
			// Check if session exists
			session_exists, err := tmuxSessionExists(project_name)
			if err != nil {
				return err
			}
			if !session_exists {
				// No session to close
				return nil
			}

			// Close the session
			_, err = shellCommand("tmux", "kill-session", "-t", project_name)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(closeCmd)
}
