package main

import (
	"fmt"
	"github.com/fatih/color"
)

// print displays a message
func print(in string, a ...interface{}) {
	if len(a) > 0 {
		in = fmt.Sprintf(in, a...)
	}

	b := color.New(color.FgHiBlue)
	fmt.Printf(" %s  %s\n", b.Sprint("◈"), in)
}
