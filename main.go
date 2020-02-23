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

	fmt.Printf("Total hours:\t\t%.2f\n", c.Sheet.TotalHours())
	fmt.Printf("Total pay:\t\t%.2f\n", c.Sheet.TotalPay())
	fmt.Printf("After tax earnings:\t$%.2f\n", c.Sheet.AfterTax())
}
