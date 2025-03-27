package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	logger *log.Logger

	workspace_layout WorkspaceLayout

	workspaces_path string
	workspaces_file string

	rootCmd = &cobra.Command{
		Use:   "workspaces",
		Short: "Quickly access your most used projects with tmux",
		Long:  `Workspaces allows you to store your project locations and quickly create named tmux sessions for them.`,
	}
)

type (
	WorkspaceLayout struct {
		Workspaces []Workspace `json:"workspaces"`
	}

	Workspace struct {
		// Name    string `json:"name"`
		Path    string `json:"path"`
		StowDir bool   `json:"stow_dir,omitempty"`
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&workspaces_path, "path", "p", "~/.local/share/workspaces", "Path to the workspaces directory")
	// TODO: properly handle directories with '~' by resolving the user home path
	workspaces_path, err := filepath.Abs(workspaces_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("TODO: fix project path resolution")
	fmt.Println(workspaces_path)
	os.Exit(1)

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
	logger = log.New(log_file, "", log.LstdFlags)

	buf, err := os.ReadFile(workspaces_file)
	if err != nil && os.IsNotExist(err) {
		fd, err := os.Create(workspaces_file)
		if err != nil {
			logger.Fatalln("Failed to read workspaces.json")
		}
		fd.Close()
	}

	if len(buf) > 0 {
		err := json.Unmarshal(buf, &workspace_layout)
		if err != nil {
			logger.Fatalln("Failed to unmarshal workspaces.json")
		}
	}

	cobra.OnFinalize(finalize)
}

func projectLinks() (map[string]string, error) {
	workspaces := workspace_layout.Workspaces
	projects := make(map[string]string)
	for _, workspace := range workspaces {
		dir_entries, err := os.ReadDir(workspace.Path)
		if err != nil {
			return nil, err
		}
		for _, entry := range dir_entries {
			if !entry.Type().IsDir() || entry.Name()[0] == '.' {
				continue
			}

			projects[entry.Name()] = filepath.Join(workspace.Path, entry.Name())

			if workspace.StowDir {
				stow_package_entries, err := os.ReadDir(filepath.Join(workspace.Path, entry.Name()))
				if err != nil {
					return nil, err
				}
				for _, stow_entry := range stow_package_entries {
					if stow_entry.Type().IsDir() && stow_entry.Name() == ".config" {
						projects[entry.Name()] = filepath.Join(workspace.Path, entry.Name(), ".config", entry.Name())
						break
					}
				}
			}
		}
	}
	return projects, nil
}

// func workspaceProjects() (map[string]string, error) {
// 	projects := make(map[string]string)
// 	for _, workspace := range workspace_layout.Workspaces {
// 		subdirs, err := os.ReadDir(workspace.Path)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// Add all non-hidden directories
// 		for _, subdir := range subdirs {
// 			if subdir.Type().IsDir() && subdir.Name()[0] != '.' {
// 				projects[subdir.Name()] = workspace.Path
// 			}
// 		}
// 	}
// 	return projects, nil
// }

func finalize() {
	str, err := json.Marshal(workspace_layout)
	if err != nil {
		logger.Fatalln("Error: Failed to marshal workspaces.")
	}

	err = os.WriteFile(workspaces_file, []byte(str), 0700)
	if err != nil {
		logger.Fatalln("Error: Failed to add changes.")
	}
}
