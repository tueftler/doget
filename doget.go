package main

import (
	"flag"
	"fmt"
	"github.com/tueftler/doget/dockerfile"
	"io"
	"os"
)

var (
	input    string
	output   string

	parser   *dockerfile.Parser
	commands = map[string]func(out io.Writer, file *dockerfile.Dockerfile) error{
		"transform": transform,
		"dump":      dump,
	}
)

func init() {
	flag.StringVar(&input, "in", "Dockerfile", "Use given dockerfile")
	flag.StringVar(&output, "out", "-", "Output (defaults to standard output)")

	parser = dockerfile.NewParser().Extend("INCLUDE", include)
}

func main() {
	flag.Parse()

	var out io.Writer
	if output == "-" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.Create(output)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
	}

	command := "transform"
	if len(flag.Args()) > 0 {
		command = flag.Args()[0]
	}

	var file dockerfile.Dockerfile
	if err := parser.ParseFile(input, &file); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	if delegate, ok := commands[command]; ok {
		if err := delegate(out, &file); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Cannot handle command `%s`\n", command)
		os.Exit(2)
	}
}
