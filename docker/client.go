package docker

import (
	"os/exec"
)

// Client represents the docker daemon
type Client interface {
	Build(args []string) error
}

// Create instantiates a new connection to the docker daemon
func Create(target string) Client {
	return &dockerCli{binary: target}
}

type dockerCli struct {
	binary string
}

func (d *dockerCli) Build(args []string) error {
	return exec.Command(d.binary, prependBuild(args)...).Run()
}

func prependBuild(args []string) []string {
	return append([]string{"build"}, args...)
}
