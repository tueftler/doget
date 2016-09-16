package build

import (
	"flag"
	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/docker"
	"github.com/tueftler/doget/dockerfile"
	"strings"
)

// BuildCommand is a thin wrapper around transform > docker build > clean
type BuildCommand struct {
	command.Command
	flags     *flag.FlagSet
	transform command.Command
	docker    docker.Client
	clean     command.Command
}

// NewCommand creates new build command instance
func NewCommand(name string, transform command.Command, clean command.Command, client docker.Client) *BuildCommand {
	return &BuildCommand{
		flags:     flag.NewFlagSet(name, flag.ExitOnError),
		transform: transform,
		docker:    client,
		clean:     clean,
	}
}

// Run performs action of build command
func (b *BuildCommand) Run(parser *dockerfile.Parser, args []string) error {
	transformArgs, dockerArgs := split(args)
	err := b.transform.Run(parser, transformArgs)
	if nil != err {
		return err
	}

	err = b.docker.Build(dockerArgs)
	if nil != err {
		return err
	}

	err = b.clean.Run(parser, []string{})
	if nil != err {
		return err
	}

	return nil
}

func split(args []string) ([]string, []string) {
	transformArgs := []string{}
	dockerArgs := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "--doget-") {
			for _, val := range strings.Split(strings.Replace(arg, "--doget", "", 1), "=") {
				transformArgs = append(transformArgs, val)
			}
		} else {
			dockerArgs = append(dockerArgs, arg)
		}
	}

	return transformArgs, dockerArgs
}
