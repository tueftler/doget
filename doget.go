package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/command/build"
	"github.com/tueftler/doget/command/clean"
	"github.com/tueftler/doget/command/dump"
	"github.com/tueftler/doget/command/transform"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/docker"
	"github.com/tueftler/doget/dockerfile"
	"github.com/tueftler/doget/provides"
	"github.com/tueftler/doget/use"
)

var (
	commands = make(map[string]command.Command)
	version  = "1.0.3"
)

func init() {
	commands["dump"] = dump.NewCommand("dump")
	commands["transform"] = transform.NewCommand("transform")
	commands["clean"] = clean.NewCommand("clean")
	commands["build"] = build.NewCommand(
		"build",
		commands["transform"],
		commands["clean"],
		docker.Create("docker"),
	)
}

func configuration(file string) (*config.Configuration, error) {
	if file == "" {
		return config.Default().Merge(config.SearchPath()...)
	} else {
		return config.Empty().MustMerge(file)
	}
}

func main() {
	var (
		cmdName    = flag.String("#1", "", "Command, one of [build, clean, dump, transform]")
		configFile = flag.String("config", "", "Configuration file to use")
	)
	flag.Parse()

	configuration, err := configuration(*configFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	*cmdName = flag.Arg(0)
	if delegate, ok := commands[*cmdName]; ok {
		parser := dockerfile.NewParser().
			Extend("USE", use.New(configuration.Repositories).Extension).
			Extend("PROVIDES", provides.Extension)

		args := flag.Args()
		if err := delegate.Run(parser, args[1:len(args)]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Cannot handle command %q\n", *cmdName)
		fmt.Printf("DoGet version %s, usage:\n", version)
		flag.PrintDefaults()
		os.Exit(2)
	}
}
