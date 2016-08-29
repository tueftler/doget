package use

import (
	"bytes"
	"fmt"
	"github.com/tueftler/doget/dockerfile"
	"io"
	"strings"
	"text/template"
)

// Statement represents a single USE statement
type Statement struct {
	Context   *Context
	Line      int
	Reference string
}

// Context represents the use context
type Context struct {
	Repositories map[string]map[string]string
}

// Origin represents the parsed components of a USE reference
type Origin struct {
	Host    string
	Vendor  string
	Name    string
	Version string
	Dir     string
	Uri     string
}

// New creates a USE instruction backed by the given repositories
func New(repositories map[string]map[string]string) *Context {
	return &Context{Repositories: repositories}
}

// String creates a string representation of an origin
func (o *Origin) String() string {
	str := o.Host + "/" + o.Vendor + "/" + o.Name
	if "" != o.Dir {
		str += "/" + o.Dir
	}
	if "" != o.Version {
		str += ":" + o.Version
	}
	return str
}

// Emit writes the USE statement
func (s *Statement) Emit(out io.Writer) {
	dockerfile.EmitInstruction(out, "USE", s.Reference)
}

// Origin parses origin from reference
func (s *Statement) Origin() (origin *Origin, err error) {
	var parsed []string
	var version, dir string

	pos := strings.LastIndex(s.Reference, ":")
	if pos == -1 {
		parsed = strings.Split(s.Reference, "/")
		version = "master"
	} else {
		parsed = strings.Split(s.Reference[0:pos], "/")
		version = s.Reference[pos+1 : len(s.Reference)]
	}

	if len(parsed) == 3 {
		dir = ""
	} else {
		dir = strings.Join(parsed[3:len(parsed)], "/")
	}

	// Compile URL
	if repository, ok := s.Context.Repositories[parsed[0]]; ok {
		origin = &Origin{Host: parsed[0], Vendor: parsed[1], Name: parsed[2], Dir: dir, Version: version}

		template, err := template.New(origin.Host).Parse(repository["url"])
		if err != nil {
			return nil, err
		}

		var uri bytes.Buffer
		if err := template.Execute(&uri, *origin); err != nil {
			return nil, err
		}

		origin.Uri = uri.String()
		return origin, nil
	} else {
		return nil, fmt.Errorf("No repository %s", parsed[0])
	}
}

// Extension func for parser
func (c *Context) Extension(file *dockerfile.Dockerfile, line int, tokens *dockerfile.Tokens) dockerfile.Statement {
	return &Statement{Context: c, Line: line, Reference: strings.Trim(tokens.NextLine(), " ")}
}
