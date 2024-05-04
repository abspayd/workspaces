package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	openCmd = &cobra.Command{
		Use:   "open [name]",
		Short: "Open a project within a workspace",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := workspaceProjects()
			if err != nil {
				return err
			}

			// == FZF ==
			fzf_query := ""
			if len(args) > 0 {
				fzf_query = args[0]
			}
			project_name, err := fzf(fzf_query, projects)
			if err != nil {
				exit_error := &exec.ExitError{}
				if errors.As(err, &exit_error) {
					// Silently return if exited
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
				// Create new session
				project_path := projects[project_name] + "/" + project_name
				_, err = shellCommand("tmux", "new-session", "-d", "-c", project_path, "-s", project_name)
				if err != nil {
					return err
				}
			}

			// Check if client is attached
			is_attached, _ := os.LookupEnv("TMUX")
			if len(is_attached) == 0 {
				err = tmuxAttach(project_name)
				if err != nil {
					return err
				}
			} else {
				// Switch to new session
				_, err = shellCommand("tmux", "switchc", "-t", project_name)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(openCmd)
}
