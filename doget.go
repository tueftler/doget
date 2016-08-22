package main

import (
	"flag"
	"fmt"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
	"os"
	"path/filepath"
)

var (
	configFile string

	parser   *dockerfile.Parser
	variants = []string{"Dockerfile.in", "Dockerfile"}
	commands = map[string]func(config *config.Configuration, args []string) error{
		"transform": transform,
		"dump":      dump,
	}
)

func init() {
	flag.StringVar(&configFile, "config", "", "Configuration file to use")

	parser = dockerfile.NewParser().Extend("USE", use)
}

func parse(input string, file *dockerfile.Dockerfile) error {
	if input == "-" {
		return parser.Parse(os.Stdin, file, "<stdin>")
	} else {
		stat, err := os.Stat(input)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			for _, name := range variants {
				variant := filepath.Join(input, name)
				_, err := os.Stat(variant)
				if err == nil {
					return parser.ParseFile(variant, file)
				}
			}
			return fmt.Errorf("Neither Dockerfile.in or Dockerfile exist in %s", input)
		} else {
			return parser.ParseFile(input, file)
		}
	}
}

func main() {
	flag.Parse()

	var configuration *config.Configuration
	var err error
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
		if err := delegate(configuration, args); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Cannot handle command `%s`\n", command)
		os.Exit(2)
	}
}
