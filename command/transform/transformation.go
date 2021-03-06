package transform

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/tueftler/doget/dockerfile"
	"github.com/tueftler/doget/provides"
	"github.com/tueftler/doget/use"
)

type Transformation struct {
	Input    string
	Output   io.Writer
	UseCache bool
}

type Provided map[string]bool

func (p Provided) add(image string) {
	p[image] = true
}

func (p Provided) contains(image string) bool {
	ok, _ := p[image]
	return ok
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
	return t.write(parser, &file, "", Provided{file.From.Image: true})
}

func prefix(paths, base string) string {
	segments := strings.Split(paths, " ")
	result := ""
	for _, segment := range segments[0 : len(segments)-1] {
		if strings.Contains(segment, "://") {
			result += segment + " "
		} else {
			result += base + segment + " "
		}
	}
	return result + segments[len(segments)-1]
}

func (t *Transformation) write(parser *dockerfile.Parser, file *dockerfile.Dockerfile, base string, provided Provided) error {
	fmt.Fprintf(os.Stderr, "Transform : %s\n", file.Source)
	for _, statement := range file.Statements {
		switch statement.(type) {
		case *provides.Statement:
			for _, image := range statement.(*provides.Statement).Images() {
				provided.add(image)
				fmt.Fprintf(os.Stderr, " ---> PROVIDES %s\n", image)
			}
			break

		case *use.Statement:
			var path string

			origin, err := statement.(*use.Statement).Origin()
			if err != nil {
				return err
			}

			path, err = fetch(origin, t.UseCache, func(transferred, total int64) {
				percentage := float64(transferred) / float64(total)
				finished := int(math.Max(percentage*float64(40), 40))
				fmt.Fprintf(
					os.Stderr,
					"\r ---> Transferring [%s%s] %.2fkB",
					strings.Repeat("#", finished),
					strings.Repeat("_", 40-finished),
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

			if !provided.contains(included.From.Image) {
				return fmt.Errorf(
					"Include %s requires %s, which was not found in provided %s",
					origin.String(),
					included.From.Image,
					reflect.ValueOf(provided).MapKeys(),
				)
			}

			dockerfile.EmitComment(t.Output, "Included from "+origin.String())
			t.write(parser, &included, filepath.ToSlash(path)+"/", provided)
			break

		// Remove "FROM"
		case *dockerfile.From:
			break

		// Prefix "ADD" paths:
		case *dockerfile.Add:
			dockerfile.EmitInstruction(t.Output, "ADD", prefix(statement.(*dockerfile.Add).Paths, base))
			break

		// Prefix "COPY" paths:
		case *dockerfile.Copy:
			dockerfile.EmitInstruction(t.Output, "COPY", prefix(statement.(*dockerfile.Copy).Paths, base))
			break

		default:
			statement.Emit(t.Output)
			break
		}
	}

	return nil
}
