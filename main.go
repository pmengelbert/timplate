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
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.SaveOutfile()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if compileFlagIsOn() {
		err = c.CompilePDF()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Printf("After tax earnings: $%.2f\n", c.Sheet.AfterTax())
}
