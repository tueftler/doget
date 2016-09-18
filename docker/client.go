package docker

import (
	"os"
	"os/exec"
)

// Client represents the docker daemon
type Client interface {
	Help() ([]byte, error)
	Build(args []string) error
}

// Create instantiates a new connection to the docker daemon
func Create(target string) Client {
	return &dockerCli{binary: target}
}

type dockerCli struct {
	binary string
}

// Help runs docker build help and returns its output
func (d *dockerCli) Help() ([]byte, error) {
	return exec.Command(d.binary, "help", "build").Output()
}

// Help executes docker build with the given arguments
func (d *dockerCli) Build(args []string) error {
	c := exec.Command(d.binary, prependBuild(args)...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		return err
	}

	if err := c.Wait(); err != nil {
		return err
	}

	return nil
}

func prependBuild(args []string) []string {
	return append([]string{"build"}, args...)
}
