package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Transformation struct {
	Output io.Writer
}

func (t *Transformation) Comment(value string) {
	fmt.Fprintf(t.Output, "# %s\n", strings.Replace(value, "\n", "\n# ", -1))
}

func (t *Transformation) Instruction(instruction, value string) {
	fmt.Fprintf(t.Output, "%s %s\n\n", instruction, strings.Replace(value, "\n", "\\\n", -1))
}

func (t *Transformation) Write(config *config.Configuration, file *dockerfile.Dockerfile, base string) error {
	for _, statement := range file.Statements {
		switch statement.(type) {
		case *Use:
			var path string
			reference := statement.(*Use).Reference

			path, err := fetch(reference, config, func(transferred, total int64) {
				percentage := float64(transferred) / float64(total)
				finished := int(math.Max(percentage*float64(20), 20))
				fmt.Fprintf(
					os.Stderr,
					"\r> Fetching %s: [%s%s] %.2fkB",
					reference,
					strings.Repeat("#", finished),
					strings.Repeat("_", 20-finished),
					float64(transferred)/float64(1024),
				)
			})
			fmt.Fprintf(os.Stderr, "\n")

			if err != nil {
				return err
			}

			var included dockerfile.Dockerfile
			if err := parse(path, &included); err != nil {
				return err
			}

			if included.From.Image != file.From.Image {
				return fmt.Errorf(
					"Include %s inherits from %s, which is incompatible with %s",
					reference,
					included.From.Image,
					file.From.Image,
				)
			}

			t.Comment("Included from " + reference)
			t.Write(config, &included, filepath.ToSlash(path)+"/")
			break

		// Retain comments
		case *dockerfile.Comment:
			t.Comment(statement.(*dockerfile.Comment).Lines)
			break

		// Builtin Docker instructions
		case *dockerfile.Maintainer:
			t.Instruction("MAINTAINER", statement.(*dockerfile.Maintainer).Name)
			break
		case *dockerfile.Run:
			t.Instruction("RUN", statement.(*dockerfile.Run).Command)
			break
		case *dockerfile.Label:
			t.Instruction("LABEL", statement.(*dockerfile.Label).Pairs)
			break
		case *dockerfile.Expose:
			t.Instruction("EXPOSE", statement.(*dockerfile.Expose).Ports)
			break
		case *dockerfile.Env:
			t.Instruction("ENV", statement.(*dockerfile.Env).Pairs)
			break
		case *dockerfile.Add:
			t.Instruction("ADD", base+statement.(*dockerfile.Add).Paths)
			break
		case *dockerfile.Copy:
			t.Instruction("COPY", base+statement.(*dockerfile.Copy).Paths)
			break
		case *dockerfile.Entrypoint:
			t.Instruction("ENTRYPOINT", statement.(*dockerfile.Entrypoint).CmdLine)
			break
		case *dockerfile.Volume:
			t.Instruction("VOLUME", statement.(*dockerfile.Volume).Names)
			break
		case *dockerfile.User:
			t.Instruction("USER", statement.(*dockerfile.User).Name)
			break
		case *dockerfile.Workdir:
			t.Instruction("WORKDIR", statement.(*dockerfile.Workdir).Path)
			break
		case *dockerfile.Arg:
			t.Instruction("ARG", statement.(*dockerfile.Arg).Name)
			break
		case *dockerfile.Onbuild:
			t.Instruction("ONBUILD", statement.(*dockerfile.Onbuild).Instruction)
			break
		case *dockerfile.Stopsignal:
			t.Instruction("STOPSIGNAL", statement.(*dockerfile.Stopsignal).Signal)
			break
		case *dockerfile.Healthcheck:
			t.Instruction("HEALTHCHECK", statement.(*dockerfile.Healthcheck).Command)
			break
		case *dockerfile.Shell:
			t.Instruction("SHELL", statement.(*dockerfile.Shell).CmdLine)
			break
		case *dockerfile.Cmd:
			t.Instruction("CMD", statement.(*dockerfile.Cmd).CmdLine)
			break
		}
	}

	return nil
}

func open(output string) (io.Writer, error) {
	if output == "-" {
		return os.Stdout, nil
	} else {
		return os.Create(output)
	}
}

func transform(config *config.Configuration, args []string) error {
	var input, output string

	flags := flag.NewFlagSet("transform", flag.ExitOnError)
	flags.StringVar(&input, "in", "Dockerfile.in", "Input. Use - for standard input")
	flags.StringVar(&output, "out", "Dockerfile", "Output. Use - for standard output")
	flags.Parse(args)

	fmt.Fprintf(os.Stderr, "> Running transform(%q -> %q) using %s\n", input, output, config.Source)

	// Parse input
	var file dockerfile.Dockerfile
	if err := parse(input, &file); err != nil {
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
	if err := transformation.Write(config, &file, ""); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Done\n\n")

	// Result
	fmt.Fprintf(out, buf.String())
	return nil
}
