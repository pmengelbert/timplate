package main

import (
	"fmt"
	"os"

	"github.com/pmengelbert/timplate/pkg/converter"
)

func main() {
	handleFlags()

	c, err := converter.DefaultConverter(infile(), outfile())
	if err != nil {
		fmt.Sprintf("%v", err)
		os.Exit(1)
	}

	err = c.SaveOutfile()
	if err != nil {
		fmt.Sprintf("%v", err)
		os.Exit(1)
	}

	if compileFlagIsOn() {
		err = c.CompilePDF()
		if err != nil {
			fmt.Sprintf("%v", err)
			os.Exit(1)
		}
	}
}
