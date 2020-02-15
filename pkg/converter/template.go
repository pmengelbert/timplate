package converter

const timesheetTemplate = `
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
