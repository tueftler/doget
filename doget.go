package main

import (
	"flag"
	"fmt"
	"github.com/tueftler/doget/dockerfile"
	"os"
)

var (
	fileName string

	functions = map[string]func(name string, file *dockerfile.Dockerfile) error{
		"run":  run,
		"dump": dump,
	}
)

func init() {
	flag.StringVar(&fileName, "file", "Dockerfile", "Use given dockerfile")
}

func run(name string, file *dockerfile.Dockerfile) error {
	return fmt.Errorf("Command `run` not yet implemented!")
}

func dump(name string, file *dockerfile.Dockerfile) error {
	fmt.Println(name)
	for _, statement := range file.Statements {
		fmt.Println(statement)
	}
	return nil
}

func main() {
	flag.Parse()

	command := "run"
	if len(flag.Args()) > 0 {
		command = flag.Args()[0]
	}

	var file dockerfile.Dockerfile
	if err := dockerfile.ParseFile(fileName, &file); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	if function, ok := functions[command]; ok {
		if err := function(fileName, &file); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Cannot handle command `%s`\n", command)
		os.Exit(2)
	}
}
