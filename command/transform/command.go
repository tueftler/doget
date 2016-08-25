package transform

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/dockerfile"
	"io"
	"os"
)

type TransformCommand struct {
	command.Command
	flags *flag.FlagSet
}

var (
	variants = []string{"Dockerfile.in", "Dockerfile"}
)

// Creates new transform command instance
func NewCommand(name string) *TransformCommand {
	return &TransformCommand{flags: flag.NewFlagSet(name, flag.ExitOnError)}
}

func open(output string) (io.Writer, error) {
	if output == "-" {
		return os.Stdout, nil
	} else {
		return os.Create(output)
	}
}

// Runs transform command
func (c *TransformCommand) Run(parser *dockerfile.Parser, args []string) error {
	var input, output string

	flags := flag.NewFlagSet("transform", flag.ExitOnError)
	flags.StringVar(&input, "in", "Dockerfile.in", "Input. Use - for standard input")
	flags.StringVar(&output, "out", "Dockerfile", "Output. Use - for standard output")
	flags.Parse(args)

	fmt.Fprintf(os.Stderr, "> Running transform(%q -> %q)\n", input, output)

	// Open output
	out, err := open(output)
	if err != nil {
		return err
	}

	// Transform
	var buf bytes.Buffer
	transformation := Transformation{Input: input, Output: &buf}
	if err := transformation.Run(parser); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Done\n\n")

	// Result
	fmt.Fprintf(out, buf.String())
	return nil
}
