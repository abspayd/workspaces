package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type (
	Workspaces struct {
		Paths map[string]Workspace `json:"paths"`
	}

	Workspace struct {
		StowDir bool `json:"stow_dir,omitempty"`
	}
)

var (
	workspaces       Workspaces
	workspaces_dirty bool

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
					log.Fatalln("Failed to create workspaces.json")
				}
				fd.Close()
			}

			if len(buf) > 0 {
				err := json.Unmarshal(buf, &workspaces)
				if err != nil {
					log.Fatalln("Failed to unmarshal workspaces.json")
				}
			}

			if workspaces.Paths == nil {
				workspaces.Paths = make(map[string]Workspace)
			}

			workspaces_dirty = false
		},
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
	workspaces := workspaces.Paths
	projects := make(map[string]string)
	for workspace, _ := range workspaces {
		dir_entries, err := os.ReadDir(workspace)
		if err != nil {
			return nil, err
		}
		for _, entry := range dir_entries {
			if !entry.Type().IsDir() || entry.Name()[0] == '.' {
				continue
			}

			projects[entry.Name()] = filepath.Join(workspace, entry.Name())

			if workspaces[workspace].StowDir {
				stow_package_entries, err := os.ReadDir(filepath.Join(workspace, entry.Name()))
				if err != nil {
					return nil, err
				}
				for _, stow_entry := range stow_package_entries {
					if stow_entry.Type().IsDir() && stow_entry.Name() == ".config" {
						projects[entry.Name()] = filepath.Join(workspace, entry.Name(), ".config", entry.Name())
						break
					}
				}
			}
		}
	}
	return projects, nil
}

func finalize() {
	if workspaces_dirty {
		str, err := json.Marshal(workspaces)
		if err != nil {
			log.Fatalln("Error: Failed to marshal workspaces.")
		}
		err = os.WriteFile(workspaces_file, []byte(str), 0700)
		if err != nil {
			log.Fatalln("Error: Failed to add changes.")
		}
	}
}
