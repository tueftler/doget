package transform

import (
	"fmt"
	"github.com/tueftler/doget/config"
	"github.com/tueftler/doget/dockerfile"
	"github.com/tueftler/doget/use"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Transformation struct {
	Input  string
	Output io.Writer
}

func parse(parser *dockerfile.Parser, input string, file *dockerfile.Dockerfile) error {
	if err := parser.ParseFile(input, file); err != nil {
		return err
	}

	if file.From == nil {
		return fmt.Errorf("File %q has no `FROM` instruction", input)
	}

	return nil
}

func load(parser *dockerfile.Parser, input string, file *dockerfile.Dockerfile) error {
	stat, err := os.Stat(input)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		for _, name := range variants {
			variant := filepath.Join(input, name)
			_, err := os.Stat(variant)
			if err == nil {
				return parse(parser, variant, file)
			}
		}
		return fmt.Errorf("Neither Dockerfile.in or Dockerfile exist in %s", input)
	} else {
		return parse(parser, input, file)
	}
}

// Run transformation
func (t *Transformation) Run(config *config.Configuration, parser *dockerfile.Parser) error {
	var file dockerfile.Dockerfile
	if err := load(parser, t.Input, &file); err != nil {
		return err
	}

	t.instruction("FROM", file.From.Image)
	return t.write(config, parser, &file, "")
}

func (t *Transformation) comment(value string) {
	fmt.Fprintf(t.Output, "# %s\n", strings.Replace(value, "\n", "\n# ", -1))
}

func (t *Transformation) instruction(instruction, value string) {
	fmt.Fprintf(t.Output, "%s %s\n\n", instruction, strings.Replace(value, "\n", "\\\n", -1))
}

func (t *Transformation) write(config *config.Configuration, parser *dockerfile.Parser, file *dockerfile.Dockerfile, base string) error {
	for _, statement := range file.Statements {
		switch statement.(type) {
		case *use.Statement:
			var path string
			reference := statement.(*use.Statement).Reference

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
			if err := load(parser, path, &included); err != nil {
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

			t.comment("Included from " + reference)
			t.write(config, parser, &included, filepath.ToSlash(path)+"/")
			break

		// Retain comments
		case *dockerfile.Comment:
			t.comment(statement.(*dockerfile.Comment).Lines)
			break

		// Builtin Docker instructions
		case *dockerfile.Maintainer:
			t.instruction("MAINTAINER", statement.(*dockerfile.Maintainer).Name)
			break
		case *dockerfile.Run:
			t.instruction("RUN", statement.(*dockerfile.Run).Command)
			break
		case *dockerfile.Label:
			t.instruction("LABEL", statement.(*dockerfile.Label).Pairs)
			break
		case *dockerfile.Expose:
			t.instruction("EXPOSE", statement.(*dockerfile.Expose).Ports)
			break
		case *dockerfile.Env:
			t.instruction("ENV", statement.(*dockerfile.Env).Pairs)
			break
		case *dockerfile.Add:
			t.instruction("ADD", base+statement.(*dockerfile.Add).Paths)
			break
		case *dockerfile.Copy:
			t.instruction("COPY", base+statement.(*dockerfile.Copy).Paths)
			break
		case *dockerfile.Entrypoint:
			t.instruction("ENTRYPOINT", statement.(*dockerfile.Entrypoint).CmdLine)
			break
		case *dockerfile.Volume:
			t.instruction("VOLUME", statement.(*dockerfile.Volume).Names)
			break
		case *dockerfile.User:
			t.instruction("USER", statement.(*dockerfile.User).Name)
			break
		case *dockerfile.Workdir:
			t.instruction("WORKDIR", statement.(*dockerfile.Workdir).Path)
			break
		case *dockerfile.Arg:
			t.instruction("ARG", statement.(*dockerfile.Arg).Name)
			break
		case *dockerfile.Onbuild:
			t.instruction("ONBUILD", statement.(*dockerfile.Onbuild).Instruction)
			break
		case *dockerfile.Stopsignal:
			t.instruction("STOPSIGNAL", statement.(*dockerfile.Stopsignal).Signal)
			break
		case *dockerfile.Healthcheck:
			t.instruction("HEALTHCHECK", statement.(*dockerfile.Healthcheck).Command)
			break
		case *dockerfile.Shell:
			t.instruction("SHELL", statement.(*dockerfile.Shell).CmdLine)
			break
		case *dockerfile.Cmd:
			t.instruction("CMD", statement.(*dockerfile.Cmd).CmdLine)
			break
		}
	}

	return nil
}
