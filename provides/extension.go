package provides

import (
	"io"
	"strings"

	"github.com/tueftler/doget/dockerfile"
)

// Statement represents a single PROVIDES statement
type Statement struct {
	Line int
	List string
}

// Emit writes the PROVIDES statement
func (s *Statement) Emit(out io.Writer) {
	dockerfile.EmitInstruction(out, "PROVIDES", s.List)
}

// Images parses the list and returns it as an array
func (s *Statement) Images() []string {
	list := strings.Split(s.List, " ")
	result := make([]string, 0)
	for _, image := range list {
		if image != "" {
			result = append(result, image)
		}
	}
	return result
}

// Extension func for parser
func Extension(file *dockerfile.Dockerfile, line int, tokens *dockerfile.Tokens) dockerfile.Statement {
	return &Statement{Line: line, List: tokens.NextLine()}
}
