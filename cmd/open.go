package cmd

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/abspayd/workspaces/internal"
	"github.com/spf13/cobra"
)

var (
	openCmd = &cobra.Command{
		Use:                   "open [name]",
		Short:                 "Open a project within a workspace",
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			projects, err := projectLinks()
			if err != nil {
				log.Println("Error: Failed to get workspace projects:", err)
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

			project_name, err := internal.Fzf(fzf_query, project_names)
			if err != nil {
				exit_error := &exec.ExitError{}
				if errors.As(err, &exit_error) {
					// Silently return if exited
					return nil
				}
				return err
			}

			// replace '.' with '_'
			project_title := strings.ReplaceAll(project_name, ".", "_")

			// == Tmux ==
			// Check if session exists
			is_attached := len(os.Getenv("TMUX")) > 0
			session_exists, err := internal.TmuxSessionExists(project_title)
			if err != nil {
				return err
			}
			if !session_exists {
				// Create new session
				project_path := projects[project_name] + "/" + project_name
				args := []string{"new-session", "-s", project_title, "-c", project_path}
				if is_attached {
					args = append(args, "-d")
				}
				shellCmd := exec.Command("tmux", args...)
				err = shellCmd.Run()
				if err != nil {
					return err
				}
			}

			log.Println("Opened project:", project_name)

			// Check if client is attached
			if !is_attached {
				err = internal.TmuxAttachSession(project_title)
				if err != nil {
					return err
				}
			} else {
				// Switch to new session
				shellCmd := exec.Command("tmux", "switch-client", "-t", project_title)
				err = shellCmd.Run()
				if err != nil {
					log.Println("Error: Unable to switch to new session:", err)
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
