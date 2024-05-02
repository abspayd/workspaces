package cmd

import (
	"errors"
	"fmt"
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

			var projects map[string]string
			projects = make(map[string]string)

			for _, workspace := range workspaces {
				entries, err := os.ReadDir(workspace)
				if err != nil {
					return err
				}

				// Add all non-hidden directories
				for _, entry := range entries {
					if entry.Type().IsDir() && entry.Name()[0] != '.' {
						projects[entry.Name()] = workspace
						// projects = append(projects, entry.Name())
					}
				}
			}

			shell, is_set := os.LookupEnv("SHELL")
			if !is_set {
				shell = "sh"
			}

			// == FZF ==
			fzf_opts := []string{"+m"}
			command := "fzf"
			for _, opt := range fzf_opts {
				command += " " + opt
			}
			shellCmd := exec.Command(shell, "-c", command)
			shellCmd.Stderr = os.Stderr
			in, _ := shellCmd.StdinPipe()
			go func() {
				for k := range projects {
					fmt.Fprintln(in, k)
				}
				in.Close()
			}()
			out, err := shellCmd.Output()
			if err != nil {
				exit := &exec.ExitError{}
				if errors.As(err, &exit) {
					return nil
				}
				return err
			}

			// Remove new line
			project_name := string(out[:len(out)-1])

			// == Tmux ==
			// Check if session exists
			filter := fmt.Sprintf("#{m:#{session_name},%s}", project_name)
			match, err := shellCommand("tmux", "ls", "-f", filter)
			if len(match) == 0 {
				// Create new session
				project_path := projects[project_name] + "/" + project_name
				_, err = shellCommand("tmux", "new-session", "-d", "-c", project_path, "-s", project_name)
				if err != nil {
					return err
				}

				fmt.Println("Created new session", project_name)
			}

			// Switch to new session
			_, err = shellCommand("tmux", "switchc", "-t", project_name)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

// Execute a shell command with n arguments
func shellCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func init() {
	rootCmd.AddCommand(openCmd)
}
