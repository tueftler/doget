package main

import (
	"flag"
	"fmt"
	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/command/dump"
	"github.com/tueftler/doget/command/transform"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
	"github.com/tueftler/doget/use"
	"os"
)

var (
	configFile string

	parser   = dockerfile.NewParser().Extend("USE", use.Extension)
	commands = map[string]command.Command{
		"dump":      dump.NewCommand("dump"),
		"transform": transform.NewCommand("transform"),
	}
)

func init() {
	flag.StringVar(&configFile, "config", "", "Configuration file to use")
}

func main() {
	flag.Parse()

	// Configfiles
	var err error
	var configuration *config.Configuration
	if configFile == "" {
		configuration, err = config.Default()
	} else {
		configuration, err = config.FromFile(configFile)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Run subcommand
	args := flag.Args()
	command := "transform"
	if len(flag.Args()) > 0 {
		command = flag.Args()[0]
		args = args[1:len(args)]
	}

	if delegate, ok := commands[command]; ok {
		if err := delegate.Run(configuration, parser, args); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Cannot handle command `%s`\n", command)
		os.Exit(2)
	}
}
