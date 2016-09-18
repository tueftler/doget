package build

import (
	"flag"
	"fmt"
	"strings"

	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/docker"
	"github.com/tueftler/doget/dockerfile"
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

func (b *BuildCommand) Usage() error {
	fmt.Println("Usage: doget build [DOGET-OPTIONS] [OPTIONS] PATH | URL | - \n")

	// Make these look like docker build --help output
	fmt.Println("  --doget-no-cache                Do not use cache")
	fmt.Println("  --doget-in                      Input (default: Dockerfile.in)")
	fmt.Println("  --doget-out                     Output (default: Dockerfile)")

	output, err := b.docker.Help()
	if err != nil {
		return err
	}

	// Only print flags usage
	for _, line := range strings.Split(string(output), "\n") {
		if strings.HasPrefix(line, "  -") {
			fmt.Println(line)
		}
	}
	fmt.Println()

	return nil
}

// Run performs action of build command
func (b *BuildCommand) Run(parser *dockerfile.Parser, args []string) error {
	if 0 == len(args) || "-help" == args[0] || "--help" == args[0] {
		return b.Usage()
	}

	transformArgs, dockerArgs := split(args)

	if err := b.transform.Run(parser, transformArgs); err != nil {
		return err
	}

	if err := b.docker.Build(dockerArgs); err != nil {
		return err
	}

	if err := b.clean.Run(parser, []string{}); err != nil {
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
