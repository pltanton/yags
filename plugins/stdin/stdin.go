package stdin

import (
	"bufio"
	"os"
)

var created bool

// Stdin type represents stdin plugin. In future it could transform to pipe
// reader plugin.
type Stdin struct {
	scanner *bufio.Scanner
	out     chan string
}

// Chan returns output of plugin
func (s Stdin) Chan() chan string { return s.out }

// NewStdin returns new instance of Stdin
func NewStdin() Stdin {
	if created {
		panic("Stdin plugin should be unique")
	}
	created = true
	return Stdin{
		scanner: bufio.NewScanner(os.Stdin),
		out:     make(chan string),
	}
}

// StartMonitor starts monitornig on stdin
func (s Stdin) StartMonitor() {
	for s.scanner.Scan() {
		s.out <- s.scanner.Text()
	}
}
