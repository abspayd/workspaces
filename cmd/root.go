package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	Logger *log.Logger

	workspaces    WorkspaceLayout
	workspace_map map[string]Workspace

	workspaces_file string

	rootCmd = &cobra.Command{
		Use:   "workspaces",
		Short: "A simple CLI workspace manager.",
		Long:  `A simple workspace manager to quickly access and modify your frequent project locations.`,
	}
)

type (
	WorkspaceLayout struct {
		Workspaces []Workspace `json:"workspaces"`
	}

	Workspace struct {
		Name string `json:"name"`
		Path string `json:"path"`
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

	workspaces_path := filepath.Join(home, ".local/share/workspaces")
	workspaces_file = filepath.Join(workspaces_path, "/workspaces.json")

	// Check if workspace directory exists
	_, err = os.Stat(workspaces_path)
	if err != nil && os.IsNotExist(err) {
		// The directory does not exit. Create the workspace directory path
		if err = os.MkdirAll(workspaces_path, 0700); err != nil {
			fmt.Println("Failed to create workspace list directory.")
			os.Exit(1)
		}
	}

	// Logging
	log_path := filepath.Join(workspaces_path, "workspaces.log")
	log_file, err := os.OpenFile(log_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create log file.")
		os.Exit(1)
	}

	Logger = log.New(log_file, "", log.LstdFlags)

	buf, err := os.ReadFile(workspaces_file)
	if err != nil && os.IsNotExist(err) {
		fd, err := os.Create(workspaces_file)
		if err != nil {
			Logger.Fatalln("Failed to read workspaces.json")
		}
		fd.Close()
	}

	if len(buf) > 0 {
		err := json.Unmarshal(buf, &workspaces)
		if err != nil {
			Logger.Fatalln("Failed to unmarshal workspaces.json")
		}
	}

	workspace_map = make(map[string]Workspace)
	for _, workspace := range workspaces.Workspaces {
		workspace_map[workspace.Name] = workspace
	}

	cobra.OnFinalize(finalize)
}

func fzf(query string, items []string) (string, error) {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
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
		for _, s := range items {
			fmt.Fprintln(in, s)
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

func tmuxSessionExists(session_name string) (bool, error) {
	filter := fmt.Sprintf("#{m:#{session_name},%s}", session_name)
	cmd := exec.Command("tmux", "ls", "-f", filter)
	match, err := cmd.Output()

	if err != nil {
		return false, err
	}
	return len(match) != 0, nil
}

func workspaceProjects() (map[string]string, error) {
	projects := make(map[string]string)
	for _, workspace := range workspaces.Workspaces {
		subdirs, err := os.ReadDir(workspace.Path)
		if err != nil {
			return nil, err
		}

		// Add all non-hidden directories
		for _, subdir := range subdirs {
			if subdir.Type().IsDir() && subdir.Name()[0] != '.' {
				projects[subdir.Name()] = workspace.Path
			}
		}

	}
	return projects, nil
}

func finalize() {
	str, err := json.Marshal(workspaces)
	if err != nil {
		Logger.Fatalln("Error: Failed to marshal workspaces.")
	}

	err = os.WriteFile(workspaces_file, []byte(str), 0700)
	if err != nil {
		Logger.Fatalln("Error: Failed to add changes.")
	}
}
