package transform

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/tueftler/doget/command"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
	"io"
	"os"
	"path/filepath"
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

func parse(parser *dockerfile.Parser, input string, file *dockerfile.Dockerfile) error {
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

func open(output string) (io.Writer, error) {
	if output == "-" {
		return os.Stdout, nil
	} else {
		return os.Create(output)
	}
}

// Runs transform command
func (c *TransformCommand) Run(config *config.Configuration, parser *dockerfile.Parser, args []string) error {
	var input, output string

	flags := flag.NewFlagSet("transform", flag.ExitOnError)
	flags.StringVar(&input, "in", "Dockerfile.in", "Input. Use - for standard input")
	flags.StringVar(&output, "out", "Dockerfile", "Output. Use - for standard output")
	flags.Parse(args)

	fmt.Fprintf(os.Stderr, "> Running transform(%q -> %q) using %s\n", input, output, config.Source)

	// Parse input
	var file dockerfile.Dockerfile
	if err := parse(parser, input, &file); err != nil {
		return err
	}

	// Open output
	out, err := open(output)
	if err != nil {
		return err
	}

	// Transform
	var buf bytes.Buffer
	transformation := Transformation{Output: &buf}
	transformation.Instruction("FROM", file.From.Image)
	if err := transformation.Write(config, parser, &file, ""); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Done\n\n")

	// Result
	fmt.Fprintf(out, buf.String())
	return nil
}
