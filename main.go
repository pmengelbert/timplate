package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pmengelbert/timplate/pkg/converter"
)

const (
	helpMessage = `Usage: timplate [OPTIONS] <infile.yaml>
OPTIONS:
	-o outfile.tex 	-- produce outfile.tex
	-c		-- compile with pdflatex (if installed)
	-h		-- show this help message
`
)

var (
	compile = flag.Bool("c", false, "compile using pdflatex (use only if installed)")
	help    = flag.Bool("h", false, "show a help message, without running the program")
	output  = flag.String("o", "timesheet.tex", "the resulting .tex file")
)

func main() {
	handleFlags()

	c, err := converter.DefaultConverter(flag.Args()[0], *output)
	if err != nil {
		fmt.Sprintf("%v", err)
		os.Exit(1)
	}

	err = c.SaveOutfile()
	if err != nil {
		fmt.Sprintf("%v", err)
		os.Exit(1)
	}

	if *compile {
		err = c.CompilePDF()
		if err != nil {
			fmt.Sprintf("%v", err)
			os.Exit(1)
		}
	}
}

func handleFlags() {
	flag.Parse()
	if *help {
		fmt.Printf(helpMessage)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		fmt.Printf("no input file specified")
		os.Exit(1)
	}
}
