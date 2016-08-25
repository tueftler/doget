package dockerfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Builtins (v1.12), see https://docs.docker.com/engine/reference/builder/
var (
	statements = map[string]func(file *Dockerfile, line int, tokens *Tokens) Statement{
		"FROM": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			file.From = &From{Line: line, Image: tokens.NextLine()}
			return file.From
		},
		"MAINTAINER": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Maintainer{Line: line, Name: tokens.NextLine()}
		},
		"RUN": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Run{Line: line, Command: tokens.NextLine()}
		},
		"CMD": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Cmd{Line: line, CmdLine: tokens.NextLine()}
		},
		"LABEL": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Label{Line: line, Pairs: tokens.NextLine()}
		},
		"EXPOSE": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Expose{Line: line, Ports: tokens.NextLine()}
		},
		"ENV": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Env{Line: line, Pairs: tokens.NextLine()}
		},
		"ADD": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Add{Line: line, Paths: tokens.NextLine()}
		},
		"COPY": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Copy{Line: line, Paths: tokens.NextLine()}
		},
		"ENTRYPOINT": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Entrypoint{Line: line, CmdLine: tokens.NextLine()}
		},
		"VOLUME": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Volume{Line: line, Names: tokens.NextLine()}
		},
		"USER": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &User{Line: line, Name: tokens.NextLine()}
		},
		"WORKDIR": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Workdir{Line: line, Path: tokens.NextLine()}
		},
		"ARG": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Arg{Line: line, Name: tokens.NextLine()}
		},
		"ONBUILD": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Onbuild{Line: line, Instruction: tokens.NextLine()}
		},
		"STOPSIGNAL": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Stopsignal{Line: line, Signal: tokens.NextLine()}
		},
		"HEALTHCHECK": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Healthcheck{Line: line, Command: tokens.NextLine()}
		},
		"SHELL": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Shell{Line: line, CmdLine: tokens.NextLine()}
		},
		"#": func(file *Dockerfile, line int, tokens *Tokens) Statement {
			return &Comment{Line: line, Lines: tokens.NextComment()}
		},
	}
)

type Parser struct {
	statements map[string]func(file *Dockerfile, line int, tokens *Tokens) Statement
	extended   bool
}

// Creates a new parser
func NewParser() *Parser {
	return &Parser{statements: statements, extended: false}
}

// Parses a dockerfile from a reader. Returns an error if
// an unknown token is encountered.
func (p *Parser) Parse(input io.Reader, file *Dockerfile, source ...string) (err error) {
	if len(source) > 0 {
		file.Source = source[0]
	} else {
		file.Source = fmt.Sprintf("%T", input)
	}

	tokens := NewTokens(input)
	for tokens.HasNext {
		token := tokens.NextToken()

		if "" == token {
			continue
		} else if statement, ok := p.statements[token]; ok {
			file.Statements = append(file.Statements, statement(file, tokens.Line, tokens))
		} else {
			return fmt.Errorf("Cannot handle token `%s` on line %d of %s", token, tokens.Line, file.Source)
		}
	}

	return nil
}

// Parses a dockerfile from a file. Returns an error if
// the file cannot be opened, is a directory or when parsing
// encounters an error
func (p *Parser) ParseFile(name string, file *Dockerfile) (err error) {
	if name == "-" {
		return Parse(os.Stdin, file, "<stdin>")
	}

	stat, err := os.Stat(name)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("The given file `%s` is a directory\n", name)
	}

	input, err := os.Open(name)
	if err != nil {
		return err
	}

	defer input.Close()
	return p.Parse(bufio.NewReader(input), file, name)
}

// Extends parser. Example:
//
//    type Include struct {
//      Line      int
//      Reference string
//    }
//
//    parser.Extend("INCLUDE", func(file *Dockerfile, line int, tokens *Tokens) Statement {
//      return &Include{Line: line, Reference: tokens.NextLine()}
//    })
//
func (p *Parser) Extend(name string, extension func(file *Dockerfile, line int, tokens *Tokens) Statement) *Parser {

	// Copy on write
	if !p.extended {
		statements := make(map[string]func(file *Dockerfile, line int, tokens *Tokens) Statement)
		for instruction, parsing := range p.statements {
			statements[instruction] = parsing
		}
		p.statements = statements
		p.extended = true
	}

	p.statements[name] = extension
	return p
}

// Convenience shortcut
func Parse(input io.Reader, file *Dockerfile, source ...string) (err error) {
	return NewParser().Parse(input, file, source...)
}

// Convenience shortcut
func ParseFile(name string, file *Dockerfile) (err error) {
	return NewParser().ParseFile(name, file)
}
