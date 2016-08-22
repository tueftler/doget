package dump

import (
	"flag"
	"fmt"
	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
)

type DumpCommand struct {
	command.Command
	flags *flag.FlagSet
}

// Creates new dump command instance
func NewCommand(name string) *DumpCommand {
	return &DumpCommand{flags: flag.NewFlagSet(name, flag.ExitOnError)}
}

// Runs dump command
func (c *DumpCommand) Run(config *config.Configuration, parser *dockerfile.Parser, args []string) error {
	c.flags.String("#1", "", "Input. Use - for standard input")
	c.flags.Parse(args)

	// Parse input
	var file dockerfile.Dockerfile
	if err := parser.ParseFile(c.flags.Arg(0), &file); err != nil {
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
