package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func Fzf(query string, items []string) (string, error) {
	cmd := exec.Command("fzf", "+m", "-q", query)
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

	// remove newline and return
	return string(out[:len(out)-1]), nil
}
