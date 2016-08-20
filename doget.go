package main

import (
	"flag"
	"fmt"
	"github.com/tueftler/doget/dockerfile"
	"os"
)

var (
	fileName string
	parser   *dockerfile.Parser
	commands = map[string]func(file *dockerfile.Dockerfile) error{
		"run":  run,
		"dump": dump,
	}
)

func init() {
	flag.StringVar(&fileName, "file", "Dockerfile", "Use given dockerfile")

	parser = dockerfile.NewParser().Extend("INCLUDE", include)
}

func run(file *dockerfile.Dockerfile) error {
	return fmt.Errorf("Command `run` not yet implemented!")
}

func dump(file *dockerfile.Dockerfile) error {
	fmt.Println(file.Source, "{")
	for _, statement := range file.Statements {
		fmt.Printf("  %T %+v\n", statement, statement)
	}
	fmt.Println("}")
	return nil
}

func main() {
	flag.Parse()

	command := "run"
	if len(flag.Args()) > 0 {
		command = flag.Args()[0]
	}

	var file dockerfile.Dockerfile
	if err := parser.ParseFile(fileName, &file); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	if delegate, ok := commands[command]; ok {
		if err := delegate(&file); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Cannot handle command `%s`\n", command)
		os.Exit(2)
	}
}
