package main

import (
	"html/template"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
)

type (
	BulletList []string

	Record struct {
		Date, Hours string
		Description BulletList
		Times       BulletList
	}

	Sheet struct {
		Name      string   `json:"name"`
		StartDate string   `json:"startDate"`
		EndDate   string   `json:"endDate"`
		Records   []Record `json:"records"`
	}
)

func main() {
	const timesheet = `
\documentclass[12pt,twoside,letterpaper]{article}
\usepackage{longtable}
\usepackage{array}
\usepackage[margin=1in]{geometry}
\begin{document}
\renewcommand{\abstractname}{Summary}
\setcounter{secnumdepth}{0}
\title{Time sheet}
\author{Peter Engelbert}
\date{<< .StartDate >> to << .EndDate >>}
\maketitle
\begin{center}
    \begin{longtable}{ m{2cm} | m{2cm} | m{2cm} || m{8cm} }
        date & times & hours & description \\
    \hline\hline
		<< range $i, $v := .Records ->>
			<< $v.Date >> & 
			<< range $a, $b := $v.Times ->> 
			<< $b >>
			<< end ->>
			& << $v.Hours >> & 
			\begin{itemize} 
			<< range $x, $y := $v.Description ->> 
				\item << $y >>
			<< end ->>
			\end{itemize}
			\\ \hline
		<< end >>
\end{longtable}
\end{center}
\end{document}
`

	x, err := ioutil.ReadFile("asdf.yaml")
	if err != nil {
		panic(err)
	}

	var s Sheet
	yaml.Unmarshal(x, &s)
	t := template.Must(template.New("timesheet").Delims("<<", ">>").Parse(timesheet))
	t.Execute(os.Stdout, s)

}
