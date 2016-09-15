package build

import (
	"flag"
	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/command/clean"
	"github.com/tueftler/doget/command/transform"
	"github.com/tueftler/doget/dockerfile"
	"os/exec"
	"strings"
)

// BuildCommand is a thin wrapper around transform > docker build > clean
type BuildCommand struct {
	command.Command
	flags *flag.FlagSet
}

// NewCommand creates new build command instance
func NewCommand(name string) *BuildCommand {
	return &BuildCommand{flags: flag.NewFlagSet(name, flag.ExitOnError)}
}

// Run performs action of build command
func (b *BuildCommand) Run(parser *dockerfile.Parser, args []string) error {
	transformArgs, dockerArgs := split(args)
	err := transform.NewCommand("transform").Run(parser, transformArgs)
	if nil != err {
		return err
	}

	err = exec.Command("docker", dockerArgs...).Run()
	if nil != err {
		return err
	}

	err = clean.NewCommand("clean").Run(parser, []string{})
	if nil != err {
		return err
	}

	return nil
}

func split(args []string) ([]string, []string) {
	transformArgs := []string{}
	dockerArgs := []string{"build"}
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
