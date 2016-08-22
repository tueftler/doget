package main

import (
	"fmt"
	"flag"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
)

func dump(config *config.Configuration, args []string) error {
	var input string

	flags := flag.NewFlagSet("transform", flag.ExitOnError)
	flags.StringVar(&input, "in", "Dockerfile.in", "Input. Use - for standard input")
 	flags.Parse(args)

 	// Parse input
	var file dockerfile.Dockerfile
	if err := parse(input, &file); err != nil {
		return err
	}

	// Dump
	fmt.Println(file.Source, "{")
	for _, statement := range file.Statements {
		fmt.Printf("  %T %+v\n", statement, statement)
	}
	fmt.Println("}")
	return nil
}
