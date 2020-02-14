package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

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

func (s Sheet) TotalHoursString() string {
	return fmt.Sprintf("%.2f", s.TotalHours())
}

func (s Sheet) TotalPay() float64 {
	return s.TotalHours() * float64(s.Rate)
}

func (s Sheet) TotalPayString() string {
	return fmt.Sprintf("%.2f", s.TotalPay())
}

const timesheet = `
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
		 & \textbf{Total:} & \textbf{<< .TotalHoursString >>} @ << .Rate ->>/hr & \textbf{Pay: \$<<- .TotalPayString >>} \\ \hline
\end{longtable}
\end{center}
\end{document}
`

func main() {
	filename := "timesheet.yaml"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error reading file")
		os.Exit(1)
	}

	var s Sheet
	yaml.Unmarshal(file, &s)
	t := template.Must(template.New("timesheet").
		Delims("<<", ">>").
		Parse(timesheet))
	buf := new(bytes.Buffer)
	t.Execute(buf, s)
	outfile := "timesheet.tex"
	ioutil.WriteFile(outfile, buf.Bytes(), 0644)
	c := exec.Command("pdflatex", outfile)
	err = c.Run()
	if err != nil {
		fmt.Println(err)
	}
	str, err := c.Output()
	fmt.Println(string(str))
}
