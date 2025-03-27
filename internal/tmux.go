package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func TmuxSessionExists(name string) (bool, error) {
	cmd := exec.Command("tmux", "ls", "-f", fmt.Sprintf("#{m:#{session_name},%s}", name))
	match, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(match) > 0, nil
}

func TmuxNewSession(name, path string) error {
	cmd := exec.Command("tmux", "new", "-s", name, "-c", path)
	err := cmd.Start()
	if err != nil {
		return err
	}

	cmd.Process.Release()

	return nil
}

func TmuxAttachSession(name string) error {
	cmd := exec.Command("tmux", "attach", "-t", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}

	cmd.Process.Release()

	return nil
}

func TmuxSwitchSession(name string) error {
	cmd := exec.Command("tmux", "switch-client", "-t", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}

	cmd.Process.Release()

	return nil
}
