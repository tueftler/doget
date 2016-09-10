package transform

import (
	"fmt"
	"github.com/tueftler/doget/dockerfile"
	"github.com/tueftler/doget/provides"
	"github.com/tueftler/doget/use"
	"io"
	"math"
	"os"
	"path/filepath"
	"reflect"
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
	return t.write(parser, &file, "", map[string]bool{file.From.Image: true})
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

func (t *Transformation) write(parser *dockerfile.Parser, file *dockerfile.Dockerfile, base string, provided map[string]bool) error {
	for _, statement := range file.Statements {
		switch statement.(type) {
		case *provides.Statement:
			for _, image := range statement.(*provides.Statement).Images() {
				provided[image] = true
				fmt.Printf("  Provides %q\n", image)
			}
			break

		case *use.Statement:
			var path string

			origin, err := statement.(*use.Statement).Origin()
			if err != nil {
				return err
			}

			path, err = fetch(origin, func(transferred, total int64) {
				percentage := float64(transferred) / float64(total)
				finished := int(math.Max(percentage*float64(20), 20))
				fmt.Fprintf(
					os.Stderr,
					"\r> Fetching %s: [%s%s] %.2fkB",
					origin.String(),
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

			if origin.As != nil {
				provided[origin.As.Image] = true
				fmt.Printf("  Provides %q via alias\n", origin.As.Image)
			}

			if ok, _ := provided[included.From.Image]; !ok {
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

		// Do not emit
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
