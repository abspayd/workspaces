package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	openCmd = &cobra.Command{
		Use:                   "open [name]",
		Short:                 "Open a project within a workspace",
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
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

			var project_names []string
			for k := range projects {
				project_names = append(project_names, k)
			}

			project_name, err := fzf(fzf_query, project_names)
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
			is_attached := len(os.Getenv("TMUX")) > 0
			session_exists, err := tmuxSessionExists(project_name)
			if err != nil {
				return err
			}
			if !session_exists {
				// Create new session
				project_path := projects[project_name] + "/" + project_name
				args := []string{"new-session", "-s", project_name, "-c", project_path}
				if is_attached {
					args = append(args, "-d")
				}
				shellCmd := exec.Command("tmux", args...)
				shellCmd.Stderr = os.Stderr
				err = shellCmd.Run()
				if err != nil {
					return err
				}
			}

			// Check if client is attached
			if !is_attached {
				err = tmuxAttach(project_name)
				if err != nil {
					return err
				}
			} else {
				// Switch to new session
				_, err = shellCommand("tmux", "switch-client", "-t", project_name)
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
