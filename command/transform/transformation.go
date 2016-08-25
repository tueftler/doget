package transform

import (
	"fmt"
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
		for _, name := range []string{"Dockerfile.in", "Dockerfile"} {
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
func (t *Transformation) Run(parser *dockerfile.Parser) error {
	var file dockerfile.Dockerfile
	if err := load(parser, t.Input, &file); err != nil {
		return err
	}

	file.From.Emit(t.Output)
	return t.write(parser, &file, "")
}

func (t *Transformation) write(parser *dockerfile.Parser, file *dockerfile.Dockerfile, base string) error {
	for _, statement := range file.Statements {
		switch statement.(type) {
		case *use.Statement:
			var path string
			reference := statement.(*use.Statement).Reference

			path, err := fetch(statement.(*use.Statement), func(transferred, total int64) {
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

			dockerfile.EmitComment(t.Output, "Included from "+reference)
			t.write(parser, &included, filepath.ToSlash(path)+"/")
			break

		// Remove "FROM"
		case *dockerfile.From:
			break

		default:
			statement.Emit(t.Output)
			break
		}
	}

	return nil
}
