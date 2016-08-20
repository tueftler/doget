package main

import (
	"flag"
	"fmt"
	"github.com/tueftler/doget/dockerfile"
	"io"
	"path/filepath"
	"os"
)

var (
	input  string
	output string

	parser   *dockerfile.Parser
	variants = []string{"Dockerfile.in", "Dockerfile"}
	commands = map[string]func(out io.Writer, file *dockerfile.Dockerfile) error{
		"transform": transform,
		"dump":      dump,
	}
)

func init() {
	flag.StringVar(&input, "in", ".", "Input. Use - for standard input")
	flag.StringVar(&output, "out", "-", "Output. Use - for standard output")

	parser = dockerfile.NewParser().Extend("INCLUDE", include)
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

func open(output string) (io.Writer, error) {
	if output == "-" {
		return os.Stdout, nil
	} else {
		return os.Create(output)
	}
}

func main() {
	flag.Parse()

	out, err := open(output)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	command := "transform"
	if len(flag.Args()) > 0 {
		command = flag.Args()[0]
	}

	var file dockerfile.Dockerfile
	if err := parse(input, &file); err != nil {
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
