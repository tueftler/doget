package dockerfile

import (
	"bufio"
	"bytes"
	"io"
)

type Tokens struct {
	reader  *bufio.Reader
	HasNext bool
	Line    int
}

var (
	eof = rune(0)
)

func NewTokens(r io.Reader) *Tokens {
	return &Tokens{reader: bufio.NewReader(r), HasNext: true, Line: 1}
}

func (t *Tokens) NextRune() rune {
	r, _, err := t.reader.ReadRune()
	if err != nil {
		t.HasNext = false
		return eof
	}

	// Handle \n (Unix), \r\n (Windows) and \r (Mac OS)
	if '\r' == r {
		bytes, err := t.reader.Peek(1)
		if err == nil && bytes[0] == '\n' {
			t.reader.Discard(1)
		}
		r = '\n'
		t.Line++
	} else if '\n' == r {
		t.Line++
	}

	return r
}

func (t *Tokens) checkComment() bool {
	peek, err := t.reader.Peek(2)
	if err == nil && peek[0] == '#' {
		if peek[1] == ' ' {
			t.reader.Discard(2)
		} else {
			t.reader.Discard(1)
		}
		return true
	} else {
		return false
	}
}

func (t *Tokens) NextToken() string {
	if t.checkComment() {
		return "#"
	}

	var buf bytes.Buffer
	for {
		if r := t.NextRune(); r == eof {
			break
		} else if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
			break
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func (t *Tokens) NextComment() string {
	var buf bytes.Buffer
	for {
		if r := t.NextRune(); r == eof {
			break
		} else if '\n' == r {

			// Comment continuation
			if t.checkComment() {
				buf.WriteRune(r)
				continue
			}

			break
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func (t *Tokens) NextLine() string {
	var buf bytes.Buffer

	continued := false
	for {
		if r := t.NextRune(); r == eof {
			break
		} else if '\\' == r {
			n := t.NextRune()

			// Line continuation
			if '\n' == n {
				buf.WriteRune(n)

				// Inline comments
				n = t.NextRune()
				continued = '#' == n
				buf.WriteRune(n)

				continue
			}

			// Regular escape
			buf.WriteRune(r)
			if n != eof {
				buf.WriteRune(n)
			}
			break
		} else if '\n' == r && !continued {
			break
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
