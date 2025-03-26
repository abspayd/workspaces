package cmd

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	closeCmd = &cobra.Command{
		Use:                   "close [name]",
		Short:                 "Close a project session",
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			shellCmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
			out, err := shellCmd.Output()
			if err != nil {
				return err
			}

			sessions := strings.Fields(string(out))

			// FZF
			fzf_query := ""
			if len(args) > 0 {
				fzf_query = args[0]
			}
			project_name, err := fzf(fzf_query, sessions)
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
			shellCmd = exec.Command("tmux", "kill-session", "-t", project_name)
			err = shellCmd.Run()
			if err != nil {
				return err
			}

			Logger.Println("Closed session:", project_name)

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(closeCmd)
}
