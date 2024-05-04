package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

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
	if err != nil && os.IsNotExist(err) {
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

func fzf(query string, item_map map[string]string) (string, error) {
	shell, is_set := os.LookupEnv("SHELL")
	if !is_set {
		shell = "sh"
	}

	fzf_opts := []string{"+m"}
	if len(query) > 0 {
		fzf_opts = append(fzf_opts, "-q", query)
	}
	command := "fzf"
	for _, opt := range fzf_opts {
		command += " " + opt
	}
	cmd := exec.Command(shell, "-c", command)
	cmd.Stderr = os.Stderr
	in, _ := cmd.StdinPipe()
	go func() {
		for k := range item_map {
			fmt.Fprintln(in, k)
		}
		in.Close()
	}()
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	project_name := string(out[:len(out)-1])

	return project_name, nil
}

func tmuxAttach(session_name string) error {
	// Not attached, so attach
	attachCmd := exec.Command("tmux", "attach", "-t", session_name)
	attachCmd.Stderr = os.Stderr
	attachCmd.Stdin = os.Stdin
	attachCmd.Stdout = os.Stdout
	err := attachCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func workspaceProjects() (map[string]string, error) {
	projects := make(map[string]string)
	for _, workspace := range workspaces {
		entries, err := os.ReadDir(workspace)
		if err != nil {
			return nil, err
		}

		// Add all non-hidden directories
		for _, entry := range entries {
			if entry.Type().IsDir() && entry.Name()[0] != '.' {
				projects[entry.Name()] = workspace
			}
		}
	}
	return projects, nil
}
func tmuxSessionExists(session_name string) (bool, error) {
	filter := fmt.Sprintf("#{m:#{session_name},%s}", session_name)
	match, err := shellCommand("tmux", "ls", "-f", filter)
	if err != nil {
		return false, err
	}
	return  len(match) != 0, nil
}

func finalize() {
	err := os.WriteFile(workspaces_list_file, []byte(strings.Join(workspaces, ",")), 0700)
	if err != nil {
		fmt.Println("Error: Failed to add changes.")
		os.Exit(1)
	}
}
