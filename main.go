package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
)

type (
	BulletList []string

	Record struct {
		Date        string     `json:"date"`
		Hours       string     `json:"hours"`
		Description BulletList `json:"description"`
		Times       BulletList `json:"times"`
	}

	Sheet struct {
		Name      string   `json:"name"`
		Rate      int      `json:"rate"`
		StartDate string   `json:"startDate"`
		EndDate   string   `json:"endDate"`
		Records   []Record `json:"records"`
	}
)

var (
	compile = flag.Bool("c", false, "compile using pdflatex (use only if installed)")
	help    = flag.Bool("h", false, "show a help message, without running the program")
	output  = flag.String("o", "timesheet.tex", "the resulting .tex file")
)

func (s Sheet) TotalHours() float64 {
	var sum float64
	for _, v := range s.Records {
		f, err := strconv.ParseFloat(v.Hours, 64)
		if err != nil {
			fmt.Println("Hours provided aren't a number")
			os.Exit(1)
		}
		sum += f
	}

	return sum
}

func (s Sheet) TotalPay() float64 {
	return s.TotalHours() * float64(s.Rate)
}

const (
	timesheet = `
\documentclass[10pt,twoside,letterpaper]{article}
\usepackage{enumitem}
\setlist[2]{nosep}
\usepackage{longtable}
\usepackage{array}
\usepackage[margin=1in]{geometry}
\begin{document}
\begin{center}
\quad \textbf{<< .Name >>},
\quad \textbf{<< .StartDate >>} to \textbf{<< .EndDate >>}
    \begin{longtable}{ m{2cm} | m{2cm} | m{2cm} || m{8cm} }
        date & times & hours & description \\
    \hline\hline
		<< range $i, $v := .Records ->>
			<< $v.Date >> & 
			<< range $a, $b := $v.Times ->> 
			<< $b >>
			<< end ->>
			& << $v.Hours >> & 
			\begin{itemize}[topsep=0pt] \itemsep0em
			<< range $x, $y := $v.Description ->> 
				\item << $y >>
			<< end ->>
			\end{itemize}
			\\ \hline
		<< end ->>
		\hline
		 & \textbf{Total:} & \textbf{<< .TotalHours | printf "%.2f"  >>} @ << .Rate ->>/hr & \textbf{Pay: \$<<- .TotalPay | printf "%.2f" >>} \\ \hline
\end{longtable}
\end{center}
\end{document}
`

	helpMessage = `Usage: timplate [OPTIONS] <infile.yaml>
OPTIONS:
	-o outfile.tex 	-- produce outfile.tex
	-c		-- compile with pdflatex (if installed)
	-h		-- show this help message
`
)

func main() {
	flag.Parse()
	if *help {
		fmt.Printf(helpMessage)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		fmt.Printf("no input file specified")
		os.Exit(1)
	}

	file, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		fmt.Println("error reading file")
		os.Exit(1)
	}

	var s Sheet
	yaml.Unmarshal(file, &s)
	s.capitalizeDescriptions()

	t := template.Must(template.New("timesheet").
		Delims("<<", ">>").
		Parse(timesheet))

	buf := new(bytes.Buffer)
	t.Execute(buf, s)
	err = ioutil.WriteFile(*output, buf.Bytes(), 0644)
	if err != nil {
		fmt.Printf("error writing file")
		os.Exit(1)
	}

	if *compile {
		compilePDF()
	}
}

func (s *Sheet) capitalizeDescriptions() {
	for i := range s.Records {
		for j, v := range s.Records[i].Description {
			a := strings.Split(v, " ")
			a[0] = strings.Title(a[0])
			str := strings.Join(a, " ")
			s.Records[i].Description[j] = str
		}
	}
}

func compilePDF() {
	c := exec.Command("pdflatex", *output)
	str, err := c.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(str))
}
