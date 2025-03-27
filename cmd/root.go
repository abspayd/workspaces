package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	workspace_layout WorkspaceLayout

	workspaces_path string
	workspaces_file string

	rootCmd = &cobra.Command{
		Use:   "workspaces",
		Short: "Quickly access your most used projects with tmux",
		Long:  `Workspaces allows you to store your project locations and quickly create named tmux sessions for them.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if strings.HasPrefix(workspaces_path, "~") {
				home_dir, err := os.UserHomeDir()
				if err != nil {
					log.Fatalln(err)
				}
				workspaces_path = filepath.Join(home_dir, workspaces_path[1:])
			}

			workspaces_path, err := filepath.Abs(workspaces_path)
			if err != nil {
				log.Fatalln(err)
			}

			workspaces_file = filepath.Join(workspaces_path, "/workspaces.json")

			// Check if workspace directory exists
			_, err = os.Stat(workspaces_path)
			if err != nil && os.IsNotExist(err) {
				// The directory does not exit. Create the workspace directory path
				if err = os.MkdirAll(workspaces_path, 0700); err != nil {
					log.Fatalln("Failed to create workspace list directory.")
				}
			}

			buf, err := os.ReadFile(workspaces_file)
			if err != nil && os.IsNotExist(err) {
				fd, err := os.Create(workspaces_file)
				if err != nil {
					log.Fatalln("Failed to read workspaces.json")
				}
				fd.Close()
			}

			if len(buf) > 0 {
				err := json.Unmarshal(buf, &workspace_layout)
				if err != nil {
					log.Fatalln("Failed to unmarshal workspaces.json")
				}
			}
		},
	}
)

type (
	WorkspaceLayout struct {
		Workspaces []Workspace `json:"workspaces"`
	}

	Workspace struct {
		Path    string `json:"path"`
		StowDir bool   `json:"stow_dir,omitempty"`
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&workspaces_path, "path", "p", "~/.local/share/workspaces", "Path to the workspaces directory")

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

func finalize() {
	str, err := json.Marshal(workspace_layout)
	if err != nil {
		log.Fatalln("Error: Failed to marshal workspaces.")
	}

	if len(str) > 0 {
		err = os.WriteFile(workspaces_file, []byte(str), 0700)
		if err != nil {
			log.Fatalln("Error: Failed to add changes.")
		}
	}
}
