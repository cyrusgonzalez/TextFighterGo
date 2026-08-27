// Package session gives every player a uniform way to talk to the game.
package session

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Session bundles one participant's input and output.
type Session struct {
	in  *bufio.Scanner
	out io.Writer
}

// New builds a Session from a reader/writer pair.
func New(r io.Reader, w io.Writer) *Session {
	return &Session{in: bufio.NewScanner(r), out: w}
}

// Printf writes formatted output to this session.
func (s *Session) Printf(format string, a ...any) {
	fmt.Fprintf(s.out, format, a...)
}

// Println writes a line to this session.
func (s *Session) Println(a ...any) {
	fmt.Fprintln(s.out, a...)
}

// Prompt writes msg, reads one line, and returns it trimmed.
// ok is false if the input stream has ended.
func (s *Session) Prompt(msg string) (line string, ok bool) {
	s.Printf("%s", msg)
	if !s.in.Scan() {
		return "", false
	}
	return strings.TrimSpace(s.in.Text()), true
}
