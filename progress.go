package main

import (
	"fmt"
	"os"
)

type Progress struct {
	noProgress bool
	counter    int
}

func NewProgress(np bool) *Progress {
	p := Progress{
		noProgress: np,
	}
	return &p
}

func (p *Progress) Increment() {
	p.counter++
}

func (p *Progress) Prints() {
	if p.noProgress {
		return
	}

	output := os.Stderr
	if p.counter%2500 == 0 {
		fmt.Fprintf(output, ".")
	}
	if p.counter%100000 == 0 {
		text := fmt.Sprintf(" %d", p.counter)
		fmt.Fprintln(output, text)
	}
}

func (p *Progress) PrintsEnd() {
	if p.noProgress {
		return
	}

	fmt.Fprintln(os.Stderr, "")
}
