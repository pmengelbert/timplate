package converter

const (
	timesheetTemplate = `
\documentclass[10pt,twoside,letterpaper]{article}
\usepackage{._enumitem}
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

	enumitem = `

%
% Copyright (C) 2003-2019 Javier Bezos http://www.texnia.com
%
% This file may be distributed and/or modified under the conditions of
% the MIT License. A version can be found at the end of this file.
%
% Repository: https://github.com/jbezos/enumitem
%
% Release
% ~~~~~~~

\NeedsTeXFormat{LaTeX2e}
\ProvidesPackage{._enumitem}[2019/06/20 v3.9 Customized lists]

% Notes
% ~~~~~
%
% The tag enit@ is used through the style
%
% To do:
% ~~~~~~
% - ref*, for adding stuff in the same fashion as label*
% - labelled descriptions (ie, label, title, body)
% - A true nextline (far from trivial and perhaps solved with
%   labelled descriptions).
% - Improved \AddEnumerateCounter
% - Compatibility with interfaces and zref-enumitem
% - "Pausing" somehow inline boxed text.
% - \@enumctr <-> \@listctr?
% - Define keys with values
% - Revise @nobreak
% - bottomsep
% - \SetEnumerateCounter - must define syntax
% - option verbose
% - collect sizes in \SetEnumitemSizes?
% - series=explicit / resume
% - package option inlinenew, to define "new" inline lists
%
% +=============================+
% |      EMULATING KEYVAL       |
% +=============================+
%
% "Thanks" to xkeyval, which uses the same macro names as keyval :-(,
% the latter has to be replicated in full here to ensure it works as
% intended. The original work is by David Carlisle, under license LPPL.
% Once the code is here, it could be optimized by adapting it to the
% specific needs of enumitem (to do).

\def\enitkv@setkeys#1#2{%
  \def\enitkv@prefix{enitkv@#1@}%
  \let\@tempc\relax
  \enitkv@do#2,\relax,}

\def\enitkv@do#1,{%
  \ifx\relax#1\empty\else
    \enitkv@split#1==\relax
    \expandafter\enitkv@do\fi}

\def\enitkv@split#1=#2=#3\relax{%
  \enitkv@@sp@def\@tempa{#1}%
  \ifx\@tempa\@empty\else
    \expandafter\let\expandafter\@tempc
      \csname\enitkv@prefix\@tempa\endcsname
    \ifx\@tempc\relax
      \enitkv@errx{\@tempa\space undefined}%
    \else
      \ifx\@empty#3\@empty
        \enitkv@default
      \else
        \enitkv@@sp@def\@tempb{#2}%
        \expandafter\@tempc\expandafter{\@tempb}\relax
      \fi
    \fi
  \fi}

\def\enitkv@default{%
  \expandafter\let\expandafter\@tempb
    \csname\enitkv@prefix\@tempa @default\endcsname
  \ifx\@tempb\relax
      \enitkv@err{No value specified for \@tempa}%
  \else
    \@tempb\relax
  \fi}

\def\enitkv@errx#1{\enit@error{#1}\@ehc}

\let\enitkv@err\enitkv@errx

\def\@tempa#1{%
  \def\enitkv@@sp@def##1##2{%
    \futurelet\enitkv@tempa\enitkv@@sp@d##2\@nil\@nil#1\@nil\relax##1}%
  \def\enitkv@@sp@d{%
    \ifx\enitkv@tempa\@sptoken
      \expandafter\enitkv@@sp@b
    \else
      \expandafter\enitkv@@sp@b\expandafter#1%
    \fi}%
  \def\enitkv@@sp@b#1##1 \@nil{\enitkv@@sp@c##1}}

\@tempa{ }

\def\enitkv@@sp@c#1\@nil#2\relax#3{\enitkv@toks@{#1}\edef#3{\the\enitkv@toks@}}

\@ifundefined{KV@toks@}
   {\newtoks\enitkv@toks@}
   {\let\enitkv@toks@\KV@toks@}

\def\enitkv@key#1#2{%
  \@ifnextchar[%
    {\enitkv@def{enumitem#1}{#2}}%
    {\@namedef{enitkv@enumitem#1@#2}####1}}

\def\enitkv@def#1#2[#3]{%
  \@namedef{enitkv@#1@#2@default\expandafter}\expandafter
    {\csname enitkv@#1@#2\endcsname{#3}}%
  \@namedef{enitkv@#1@#2}##1}

% This ends the code copied from keyval (under LPPL).

% +=============================+
% |        DEFINITIONS          |
% +=============================+
%
% (1) The package uses a token register very often. To be on the
%     safe side, instead of \toks@, etc., a new one is declared.
% (2) \enit@inbox is the box storing the items in boxed inline
%     lists.
% (3) \enit@outerparindent is used to save the outer parindent
%     so that it can be used in the key parindent
% (4) \enit@type has three values: 0 = enum, 1 = item, 2 = desc.
% (5) \enit@calc stores which dimen is to be computed:
%     0=labelindent, 1=labelwidth, 2=labelsep, 3=leftmargin,
%     4=itemindent
% (6) \enit@resuming has four values: 0 = none, 1 = series,
%     2 = resume* series (computed in group enumitem-resume),
%     3 = resume* list (ie, with no value).

\chardef  \enit@iv=4
\newlength\labelindent
\newdimen \enit@outerparindent
\newtoks  \enit@toks
\newbox   \enit@inbox

\newif\ifenit@boxmode
\newif\ifenit@sepfrommargin
\newif\ifenit@lblfrommargin
\newif\ifenit@calcwidest
\newif\ifenit@nextline
\newif\ifenit@boxdesc

% An alias (calc-savvy):

\let\c@enit@cnt\@tempcnta

\def\enit@meaning{\expandafter\strip@prefix\meaning}
\def\enit@noexcs#1{\expandafter\noexpand\csname#1\endcsname}

\long\def\enit@afterelse#1\else#2\fi{\fi#1}
\long\def\enit@afterfi#1\fi{\fi#1}
\def\enit@ifunset#1{%
  \expandafter\ifx\csname#1\endcsname\relax
    \expandafter\@firstoftwo
  \else
    \expandafter\@secondoftwo
  \fi}
\enit@ifunset{ifcsname}%
  {}%
  {\def\enit@ifunset#1{%
     \ifcsname#1\endcsname
       \expandafter\ifx\csname#1\endcsname\relax
         \enit@afterelse\expandafter\@firstoftwo
       \else
         \enit@afterfi\expandafter\@secondoftwo
       \fi
     \else
       \expandafter\@firstoftwo
     \fi}}

% Miscelaneous errors
% ===================

\def\enit@error{\PackageError{enumitem}}

\def\enit@checkerror#1#2{%
  \enit@error{Unknown value '#2' for key '#1'}%
      {See the manual for valid values}}

\def\enit@itemerror{%
  \enit@error{Misplaced \string\item}%
      {Either there is some text before the first\MessageBreak
       item or the last item has no text}}

\def\enit@noserieserror#1{%
  \enit@error{Series '#1' not started}%
      {You are trying to continue a series\MessageBreak
       which has not been started with 'series'}}

\def\enit@checkseries#1{%
  \ifcase\enit@resuming
    \enit@error{Misplaced key '#1'}%
      {'series' and 'resume*' must be used\MessageBreak
       in the optional argument of lists}%
  \fi}

\def\enit@checkseries@m{%
  \ifcase\enit@resuming\else
    \enit@error{Uncompatible series settings}%
      {'series' and 'resume*' must not be used\MessageBreak
       at the same time}%
  \fi}

\let\enit@toodeep\@toodeep

\def\@toodeep{%
  \ifnum\@listdepth>\enit@listdepth\relax
    \enit@toodeep
  \else
    \count@\@listdepth
    \global\advance\@listdepth\@ne
    \enit@ifunset{@list\romannumeral\the\@listdepth}%
      {\expandafter\let
         \csname @list\romannumeral\the\@listdepth\expandafter\endcsname
         \csname @list\romannumeral\the\count@\endcsname}{}%
  \fi}

% +=============================+
% |            KEYS             |
% +=============================+
%
% Including code executed by keys.
%
% There are 2 keyval groups: enumitem, and enumitem-delayed.
% The latter is used to make sure a prioritary key is the latest one;
% eg, ref, so that the ref format set by label is overriden. So, when
% this key is found in enumitem, nothing is done, except the key/value
% is moved to enumitem-delayed.
%
% A further group (enumitem-resume) catches resume* and series in
% optional arguments in lists.
%
% Vertical spacing
% ================

\enitkv@key{}{topsep}{%
  \enit@setlength\topsep{#1}}

\enitkv@key{}{itemsep}{%
  \enit@setlength\itemsep{#1}}

\enitkv@key{}{parsep}{%
  \enit@setlength\parsep{#1}}

\enitkv@key{}{partopsep}{%
  \enit@setlength\partopsep{#1}}

% Horizontal spacing
% ==================
%
% There are 3 cases: *, ! and a value. The latter also
% cancels widest with the sequence key=* ... key=value
% \string is used, just in case some package changes the
% catcodes.

\def\enit@calcset#1#2#3{%
  \if\string*\string#3%
    \enit@calcwidesttrue
    \let\enit@calc#2%
  \else\if\string!\string#3%
    \enit@calcwidestfalse
    \let\enit@calc#2%
  \else
    \ifnum\enit@calc=#2%
      \enit@calcwidestfalse
      \let\enit@calc\z@
    \fi
    \enit@setlength#1{#3}%
  \fi\fi}

\def\enitkv@enumitem@widest#1{%
  \ifcase\enit@type  % enumerate
    \expandafter\let\csname enit@cw@\@enumctr\endcsname\relax
    \@namedef{enit@widest@\@enumctr}##1{\enit@format{#1}}%
  \else              % itemize / description
    \def\enit@widest@{\enit@format{#1}}%
  \fi}

\def\enitkv@enumitem@widest@default{%
  \expandafter\let\csname enit@cw@\@enumctr\endcsname\relax
  \expandafter\let\csname enit@widest@\@enumctr\endcsname\relax}

\enitkv@key{}{widest*}{%
  \setcounter{enit@cnt}{#1}%
  \expandafter\edef\csname enit@cw@\@enumctr\endcsname
    {\the\c@enit@cnt}%
  \expandafter\edef\csname enit@widest@\@enumctr\endcsname##1%
    {##1{\the\c@enit@cnt}}}

\enitkv@key{}{labelindent*}{%
  \enit@lblfrommargintrue
  \ifnum\enit@calc=\z@
    \enit@calcwidestfalse
  \fi
  \enit@setlength\labelindent{#1}%
  \advance\labelindent\leftmargin}

\enitkv@key{}{labelindent}{%
  \enit@lblfrommarginfalse
  \enit@calcset\labelindent\z@{#1}}

\enitkv@key{}{labelwidth}{%
  \enit@calcset\labelwidth\@ne{#1}}

\enitkv@key{}{leftmargin}{%
  \edef\enit@c{\the\leftmargin}%
  \enit@calcset\leftmargin\thr@@{#1}%
  \ifenit@lblfrommargin
    \advance\labelindent-\enit@c\relax
    \advance\labelindent\leftmargin
  \fi}

\enitkv@key{}{itemindent}{%
  \edef\enit@c{\the\itemindent}%
  \enit@calcset\itemindent\enit@iv{#1}%
  \ifenit@sepfrommargin
    \advance\labelsep-\enit@c\relax
    \advance\labelsep\itemindent
  \fi}

\enitkv@key{}{listparindent}{%
  \enit@setlength\listparindent{#1}}

\enitkv@key{}{rightmargin}{%
  \enit@setlength\rightmargin{#1}}

% labelsep, from itemindent; labelsep*, from leftmargin

\enitkv@key{}{labelsep*}{%
  \enit@sepfrommargintrue
  \ifnum\enit@calc=\tw@
    \enit@calcwidestfalse
    \let\enit@calc\z@
  \fi
  \enit@setlength\labelsep{#1}%
  \advance\labelsep\itemindent}

\enitkv@key{}{labelsep}{%
  \enit@sepfrommarginfalse
  \enit@calcset\labelsep\tw@{#1}}

\enitkv@key{}{left}{%
  \enit@setleft#1..\@empty..\@@}

\def\enit@setleft#1..#2..#3\@@{%
  \enit@setlength\labelindent{#1}%
  \edef\enit@a{#3}%
  \ifx\enit@a\@empty
    \enit@calcset\leftmargin\thr@@*%
  \else
    \enit@setlength\leftmargin{#2}%
    \enit@calcset\labelsep\tw@*%
  \fi}

% Series, resume and start
% ========================

\enitkv@key{-resume}{series}{%
  \enit@checkseries@m
  \let\enit@resuming\@ne  %%% TODO - default check also \Set..Key
  \ifcase\enit@seriesopt
    \enit@ifunset{enitkv@enumitem@#1}{}%
      {\enit@error
        {Invalid series name '#1'}%
        {Do not name a series with an existing key}}%
  \else  % series=override
    \global\@namedef{enitkv@enumitem@#1}%    with value
      {\enit@error
        {Key '#1' has been overriden by a series}%
        {Change the series name and/or deactivate series=override}}%
    \global\@namedef{enitkv@enumitem@#1@default}{}%
  \fi
  \def\enit@series{#1}}

\enitkv@key{}{series}{%
  \enit@checkseries{series}}

\def\enitkv@enumitem@resume#1{%
  \edef\enit@series{#1}%
  \@nameuse{enit@resume@series@#1}\relax}

\def\enitkv@enumitem@resume@default{%
  \@nameuse{enit@resume@\@currenvir}\relax}

\@namedef{enitkv@enumitem-resume@resume*}#1{%
  \enit@checkseries@m
  \let\enit@resuming\tw@
  \edef\enit@series{#1}%
  \enit@ifunset{enit@resumekeys@series@#1}%
    {\enit@noserieserror{#1}}%
    {\expandafter\let\expandafter\enit@resumekeys
         \csname enit@resumekeys@series@#1\endcsname}}

\@namedef{enitkv@enumitem-resume@resume*@default}{%
  \let\enit@resuming\thr@@
  \expandafter\let\expandafter\enit@resumekeys
    \csname enit@resumekeys@\@currenvir\endcsname
  \@nameuse{enit@resume@\@currenvir}\relax}

\enitkv@key{}{resume*}[]{%
  \enit@checkseries{resume*}}

\newcommand\restartlist[1]{%
  \enit@ifunset{end#1}%
    {\enit@error{Undefined list '#1'}%
      {No list has been defined with that name.}}%
    {\expandafter\let
     \csname enit@resume@#1\endcsname\@empty}}

\enitkv@key{}{start}[\@ne]{%
  \setcounter{\@listctr}{#1}%
  \advance\@nameuse{c@\@listctr}\m@ne}

% Penalties
% =========

\enitkv@key{}{beginpenalty}{%
  \@beginparpenalty#1\relax}

\enitkv@key{}{midpenalty}{%
  \@itempenalty#1\relax}

\enitkv@key{}{endpenalty}{%
  \@endparpenalty#1\relax}

% Font/Format
% ===========

\enitkv@key{}{format}{%
  \def\enit@format{#1}}

\enitkv@key{}{font}{%
  \def\enit@format{#1}}

% Description styles
% ==================

\enitkv@key{}{style}[normal]{%
  \enit@ifunset{enit@style@#1}%
    {\enit@checkerror{style}{#1}}%
    {\enit@nextlinefalse
     \enit@boxdescfalse
     \@nameuse{enit@style@#1}%
     \edef\enit@descstyle{\enit@noexcs{enit@#1style}}}}

\def\enit@style@standard{%
  \enit@boxdesctrue
  \enit@calcset\itemindent\enit@iv!}

\let\enit@style@normal\enit@style@standard

\def\enit@style@unboxed{%
  \enit@calcset\itemindent\enit@iv!}

\def\enit@style@sameline{%
  \enit@calcset\labelwidth\@ne!}

\def\enit@style@multiline{%
  \enit@align@parleft
  \enit@calcset\labelwidth\@ne!}

\def\enit@style@nextline{%
  \enit@nextlinetrue
  \enit@calcset\labelwidth\@ne!}

% Labels and refs
% ===============

% Aligment
% --------

\enitkv@key{}{align}{%
  \enit@ifunset{enit@align@#1}%
    {\enit@checkerror{align}{#1}}%
    {\csname enit@align@#1\endcsname}}

% \nobreak for unboxed label with color. See below.

\newcommand\SetLabelAlign[2]{%
  \enit@toks{#2}%
  \expandafter\edef\csname enit@align@#1\endcsname
    {\def\noexpand\enit@align####1{\nobreak\the\enit@toks}}}

\def\enit@align@right{%
  \def\enit@align##1{\nobreak\hss\llap{##1}}}

\def\enit@align@left{%
  \def\enit@align##1{\nobreak##1\hfil}}

\def\enit@align@parleft{%
  \def\enit@align##1{%
    \nobreak
    \strut\smash{\parbox[t]\labelwidth{\raggedright##1}}}}

% \enit@ref has three possible definitions:
% (1) \relax, if there is neither label nor ref (ie, use
%   LaTeX settings).
% (2) set ref to @itemlabel, if there is label but not ref
% (3) set ref to ref, if there is ref (with or without label)

\enitkv@key{}{label}{%
  \expandafter\def\@itemlabel{#1}%
  \def\enit@ref{\expandafter\enit@reflabel\@itemlabel\z@}}

\enitkv@key{}{label*}{%
  \ifnum\enit@depth=\@ne
    \expandafter\def\@itemlabel{#1}%
  \else % no level 0
    \advance\enit@depth\m@ne
    \enit@toks{#1}%
    \expandafter\edef\@itemlabel{%
      \enit@noexcs{label\enit@prevlabel}%
      \the\enit@toks}%
    \advance\enit@depth\@ne
  \fi
  \def\enit@ref{\expandafter\enit@reflabel\@itemlabel\z@}}

% ref is set by label, except if there is an explicit ref in the same
% hierarchy level. Explicit refs above the current hierarchy level are
% overriden by label (besides ref), too. Since an explicit ref has
% preference, it's delayed.

\enitkv@key{}{ref}{%
  \g@addto@macro\enit@delayedkeys{,ref=#1}}

\enitkv@key{-delayed}{ref}{%
  \def\enit@ref{\enit@reflabel{#1}\@ne}}

% #2=0 don't "normalize" (ie, already normalized)
%   =1 "normalize" (in key ref)
% Used thru \enit@ref

\def\enit@reflabel#1#2{%
  \ifnum\enit@depth=\@ne\else % no level 0
    \advance\enit@depth\@ne
    \@namedef{p@\@enumctr}{}% Don't accumulate labels
    \advance\enit@depth\m@ne
  \fi
  \ifcase#2%
    \@namedef{the\@enumctr}{{#1}}%
  \else
    \enit@normlabel{\csname the\@enumctr\endcsname}{#1}%
  \fi}

% \xxx* in counters (refstar) and widest (calcdef)
% ------------------------------------------------
% \enit@labellist contains a list of
% \enit@elt{widest}\count\@count\enit@sc@@count
% \enit@elt is either \enit@getwidth or \enit@refstar, defined
% below
% The current implementation is sub-optimal -- labels are stored in
% labellist, counters defined again when processing labels, and
% modifying it is almost impossible.

\let\enit@labellist\@empty

\newcommand\AddEnumerateCounter{%
  \@ifstar\enit@addcounter@s\enit@addcounter}

\def\enit@addcounter#1#2#3{%
  \enit@toks\expandafter{%
    \enit@labellist
    \enit@elt{#3}}%
  \edef\enit@labellist{%
    \the\enit@toks
    \enit@noexcs{\expandafter\@gobble\string#1}%
    \enit@noexcs{\expandafter\@gobble\string#2}%
    \enit@noexcs{enit@sc@\expandafter\@gobble\string#2}}}

\def\enit@addcounter@s#1#2#3{%
  \enit@addcounter{#1}{#2}%
    {\@nameuse{enit@sc@\expandafter\@gobble\string#2}{#3}}}

% The 5 basic counters:

\AddEnumerateCounter\arabic\@arabic{0}
\AddEnumerateCounter\alph\@alph{m}
\AddEnumerateCounter\Alph\@Alph{M}
\AddEnumerateCounter\roman\@roman{viii}
\AddEnumerateCounter\Roman\@Roman{VIII}

% Inline lists
% ============
%
% Labels
% ------

\enitkv@key{}{itemjoin}{%
  \def\enit@itemjoin{#1}}

\enitkv@key{}{itemjoin*}{%
  \def\enit@itemjoin@s{#1}}

\enitkv@key{}{afterlabel}{%
  \def\enit@afterlabel{#1}}

% Mode
% ----

\enitkv@key{}{mode}{%
  \enit@ifunset{enit@mode#1}%
    {\enit@checkerror{mode}{#1}}%
    {\csname enit@mode#1\endcsname}}

\let\enit@modeboxed\enit@boxmodetrue
\let\enit@modeunboxed\enit@boxmodefalse

% Short Labels
% ============

\let\enit@marklist\@empty

% shorthand, expansion:

\newcommand\SetEnumerateShortLabel[2]{%
  \let\enit@a\@empty
  \def\enit@elt##1##2{%
    \def\enit@b{#1}\def\enit@c{##1}%
    \ifx\enit@b\enit@c\else
      \expandafter\def\expandafter\enit@a\expandafter{%
        \enit@a
        \enit@elt{##1}{##2}}%
    \fi}%
  \enit@marklist
  \expandafter\def\expandafter\enit@a\expandafter{%
    \enit@a
    \enit@elt{#1}{#2}}%
  \let\enit@marklist\enit@a}

\SetEnumerateShortLabel{a}{\alph*}
\SetEnumerateShortLabel{A}{\Alph*}
\SetEnumerateShortLabel{i}{\roman*}
\SetEnumerateShortLabel{I}{\Roman*}
\SetEnumerateShortLabel{1}{\arabic*}

% This is called \enit@first one,two,three,\@nil\@@nil. If there
% are just one element #2 is \@nil, otherwise we have to remove
% the trailing ,\@nil with enit@first@x
% Called with the keys in \enit@c
% Returns enit@toks

\def\enit@first#1,#2\@@nil{%
  \in@{=}{#1}% Quick test, if contains =, it's key=value
  \ifin@\else
    \enitkv@@sp@def\enit@a{#1}%
    \enit@ifunset{enitkv@enumitem@\enit@meaning\enit@a}%
      {\ifnum\enit@type=\z@
         \def\enit@elt{\enit@replace\enit@a}%
         \enit@marklist % Returns \enit@toks
       \else
         \enit@toks{#1}%
       \fi
       \ifx\@nil#2%
         \ifx,#1,\else
           \edef\enit@c{label=\the\enit@toks}%
         \fi
       \else
         \@temptokena\expandafter{\enit@first@x#2}%
         \edef\enit@c{label=\the\enit@toks,\the\@temptokena}%
       \fi}%
     {}%
  \fi
  \enit@toks\expandafter{\enit@c}}

\def\enit@first@x#1,\@nil{#1}

\def\enit@replace#1#2#3{%
  \enit@toks{}%
  \def\enit@b##1#2##2\@@nil{%
    \ifx\@nil##2%
      \addto@hook\enit@toks{##1}%
    \else
      \edef\enit@a{\the\enit@toks}%
      \ifx\enit@a\@empty\else
        \enit@error{Extra short label ignored}%
           {There are more than one short label}%
      \fi
      \addto@hook\enit@toks{##1#3}%
      \enit@b##2\@@nil
    \fi}%
  \expandafter\enit@b#1#2\@nil\@@nil
  \edef#1{\the\enit@toks}}

% Pre and post code
% =================

\enitkv@key{}{before}{%
  \def\enit@before{#1}}

\enitkv@key{}{before*}{%
  \expandafter\def\expandafter\enit@before\expandafter
    {\enit@before#1}}

\enitkv@key{}{after}{%
  \def\enit@after{#1}}

\enitkv@key{}{after*}{%
  \expandafter\def\expandafter\enit@after\expandafter
    {\enit@after#1}}

\enitkv@key{}{first}{%
  \def\enit@keyfirst{#1}}

\enitkv@key{}{first*}{%
  \expandafter\def\expandafter\enit@keyfirst\expandafter
    {\enit@keyfirst#1}}

% Miscelaneous keys
% ================

\enitkv@key{}{nolistsep}[true]{%
  \partopsep=\z@skip
  \topsep=\z@ plus .1pt
  \itemsep=\z@skip
  \parsep=\z@skip}

\enitkv@key{}{nosep}[true]{%
  \partopsep=\z@skip
  \topsep=\z@skip
  \itemsep=\z@skip
  \parsep=\z@skip}

\enitkv@key{}{noitemsep}[true]{%
  \itemsep=\z@skip
  \parsep=\z@skip}

\enitkv@key{}{wide}[\parindent]{%
  \enit@align@left
  \leftmargin\z@
  \labelwidth\z@
  \enit@setlength\labelindent{#1}%
  \listparindent\labelindent
  \enit@calcset\itemindent\enit@iv!}

% The following is deprecated in favour of wide:

\enitkv@key{}{fullwidth}[true]{%
  \leftmargin\z@
  \labelwidth\z@
  \def\enit@align##1{\hskip\labelsep##1}}

% "Abstract" layer
% ================
%
% Named values
% ------------

\newcommand\SetEnumitemValue[2]{% Implicit #3
  \enit@ifunset{enit@enitkv@#1}%
    {\enit@ifunset{enitkv@enumitem@#1}%
       {\enit@error{Wrong key '#1' in \string\SetEnumitemValue}%
          {Perhaps you have misspelled it}}{}%
     \expandafter\let\csname enit@enitkv@#1\expandafter\endcsname
       \csname enitkv@enumitem@#1\endcsname}{}%
  \@namedef{enitkv@enumitem@#1}##1{%
    \def\enit@a{##1}%
    \enit@ifunset{enit@enitkv@#1@\enit@meaning\enit@a}%
      {\@nameuse{enit@enitkv@#1}{##1}}%
      {\@nameuse{enit@enitkv@#1\expandafter\expandafter\expandafter}%
         \expandafter\expandafter\expandafter
         {\csname enit@enitkv@#1@##1\endcsname}}{}}%
  \@namedef{enit@enitkv@#1@#2}}

% Defining keys
% -------------

\newcommand\SetEnumitemKey[2]{%
  \enit@ifunset{enitkv@enumitem@#1}%
    {\enitkv@key{}{#1}[]{\enitkv@setkeys{enumitem}{#2}}}%
    {\enit@error{Duplicated key '#1' in \string\SetEnumitemKey}%
       {There already exists a key with that name}}}

% +=============================+
% |       PROCESSING KEYS       |
% +=============================+
%
% Set keys
% ========
%
% Default definition. Modified below with package option 'sizes'.

\def\enit@setkeys#1{%
  \enit@ifunset{enit@@#1}{}%
    {\expandafter\expandafter\expandafter
     \enit@setkeys@i\csname enit@@#1\endcsname\@@}}

% The following is used directly in resumeset:

\def\enit@setkeys@i#1\@@{%
  \let\enit@delayedkeys\@empty
  \enit@shl{#1}% is either \enit@toks or returns it
  \expandafter\enit@setkeys@ii\the\enit@toks\@@}

\def\enit@setkeys@ii#1\@@{%
  \enitkv@setkeys{enumitem}{#1}%
  \enit@toks\expandafter{\enit@delayedkeys}%
  \edef\enit@a{%
    \noexpand\enitkv@setkeys{enumitem-delayed}{\the\enit@toks}}%
  \enit@a}

% Handling * and ! values
% =======================
%
% \@gobbletwo removes \c from \c@counter.

\def\enit@getwidth#1#2#3#4{%
  \let#4#3%
  \def#3##1{%
    \enit@ifunset{enit@widest\expandafter\@gobbletwo\string##1}% if no widest=key
      {#1}%
      {\csname enit@widest\expandafter\@gobbletwo\string##1\endcsname{#4}}}}

\def\enit@valueerror#1{\z@ % if after an assignment, but doesn't catch \ifnum
   \enit@error{No default \string\value\space for '#1'}%
     {You can provide one with widest*}}%

\let\enit@values\@empty

\def\enit@calcwidth{%
  \ifenit@calcwidest
    \ifcase\enit@type   % ie, enum
      \enit@ifunset{enit@cw@\@enumctr}%
        {\@namedef{enit@cv@\@enumctr}{\enit@valueerror\@enumctr}}%
        {\edef\enit@values{%
           \enit@values
           \@nameuse{c@\@enumctr}\@nameuse{enit@cw@\@enumctr}\relax}%
         \expandafter
         \edef\csname enit@cv@\@enumctr\endcsname
           {\@nameuse{c@\@enumctr}}}%
      \begingroup
        \enit@values
        \def\value##1{\csname enit@cv@##1\endcsname}%
        \let\enit@elt\enit@getwidth
        \enit@labellist
        \settowidth\labelwidth{\@itemlabel}%
        \xdef\enit@a{\labelwidth\the\labelwidth\relax}%
      \endgroup
      \enit@a
    \or                 % ie, item
      \ifx\enit@widest@\relax
        \settowidth\labelwidth{\@itemlabel}%
      \else
        \settowidth\labelwidth{\enit@widest@}%
      \fi
    \else               % ie, desc
      \ifx\enit@widest@\relax
        \settowidth\labelwidth{\@itemlabel}%
      \else
        \settowidth\labelwidth{\makelabel{\enit@widest@}}%
      \fi
      \advance\labelwidth-\labelsep
    \fi
  \fi
  \advance\dimen@-\labelwidth}

\def\enit@calcleft{%
  \dimen@\leftmargin
  \advance\dimen@\itemindent
  \advance\dimen@-\labelsep
  \advance\dimen@-\labelindent
  \ifcase\enit@calc % = 0 = labelindent
    \enit@calcwidth
    \advance\labelindent\dimen@
  \or % = 1 = labelwidth, so no \enit@calcwidth
    \labelwidth\dimen@
  \or % = 2 = labelsep
    \enit@calcwidth
    \advance\labelsep\dimen@
  \or % = 3 = leftmargin
    \enit@calcwidth
    \advance\leftmargin-\dimen@
  \or % = 4 =itemindent
    \enit@calcwidth
    \advance\itemindent-\dimen@
  \fi}
  
\def\enit@negwidth{%
  \ifdim\labelwidth<\z@
    \PackageWarning{enumitem}%
       {Negative labelwidth. This does not make much\MessageBreak
        sense,}%
  \fi}

% "Normalizing" labels
% ====================
%
% Replaces \counter* by \counter{level} (those in \enit@labellist).
%
% #1 is either \csname...\endcsmame or the container \@itemlabel --
% hence \expandafter

\def\enit@refstar@i#1#2{%
  \if*#2\@empty
    \noexpand#1{\@enumctr}%
  \else
    \noexpand#1{#2}%
  \fi}%

\def\enit@refstar#1#2#3#4{%
  \def#2{\enit@refstar@i#2}%
  \def#3{\enit@refstar@i#3}}

\def\enit@normlabel#1#2{%
  \begingroup
    \def\value{\enit@refstar@i\value}%
    \let\enit@elt\enit@refstar
    \enit@labellist
    \protected@xdef\enit@a{{#2}}% Added braces as \ref is in the
  \endgroup
  \expandafter\let#1\enit@a}                    % global scope.

% Preliminary settings and default values
% =======================================

\def\enit@prelist#1#2#3{%
  \let\enit@type#1%
  \def\enit@depth{#2}%
  \edef\enit@prevlabel{#3\romannumeral#2}%
  \advance#2\@ne}

\newcount\enit@count@id

\def\enit@tagid{%
  \global\advance\enit@count@id\@ne
  \edef\EnumitemId{\number\enit@count@id}}

\def\enit@preset#1#2#3{%
   \enit@tagid
   \enit@sepfrommarginfalse
   \enit@calcwidestfalse
   \let\enit@widest@\relax
   \let\enit@resuming\z@
   \let\enit@series\relax
   \enit@boxmodetrue
   \def\enit@itemjoin{ }%
   \let\enit@itemjoin@s\relax
   \let\enit@afterlabel\nobreakspace
   \let\enit@before\@empty
   \let\enit@after\@empty
   \let\enit@keyfirst\@empty
   \let\enit@format\@firstofone % and NOT empty
   \let\enit@ref\relax
   \labelindent\z@skip
   \ifnum\@listdepth=\@ne
     \enit@outerparindent\parindent
   \else
     \parindent\enit@outerparindent
   \fi
   \enit@setkeys{list}%
   \enit@setkeys{list\romannumeral\@listdepth}%
   \enit@setkeys{#1}%
   \enit@setkeys{#1\romannumeral#2}%
   \enit@setresume{#3}}

% keyval "error" in enumitem-resume: all undefined keys (ie, all
% except resume*) are ignored, but <series> is treated like
% resume*=<series>

\def\enitkv@err@a#1{%
   \enit@ifunset{enit@resumekeys@series@\@tempa}{}%
     {\@nameuse{enitkv@enumitem-resume@resume*\expandafter}%
        \expandafter{\@tempa}}}

% keyval "error" in the optional argument: all undefined keys are
% passed to the keyval error, but <series> is ignored (already
% processed in enumitem-resume)

\def\enitkv@err@b#1{%
   \enit@ifunset{enit@resumekeys@series@\@tempa}%
     {\enit@savekverr{#1}}%
     {}}

% Process keys in optional argument:

\def\enit@setresume#1{%
  \enit@shl{#1}% Returns enit@toks
  \edef\enit@savekeys{\the\enit@toks}%
  \let\enit@savekverr\enitkv@errx
  \let\enitkv@errx\enitkv@err@a
  \edef\enit@b{%
    \noexpand\enitkv@setkeys{enumitem-resume}{\the\enit@toks}}%
  \enit@b
  \let\enitkv@errx\enitkv@err@b
  \ifcase\enit@resuming\or\or % = 2, resume* series
    \expandafter
    \enit@setkeys@i\enit@resumekeys,resume=\enit@series\@@
  \or % = 3
    \expandafter
    \enit@setkeys@i\enit@resumekeys,resume\@@
  \fi
  \expandafter\enit@setkeys@i\enit@savekeys\@@
  \let\enitkv@errx\enit@savekverr}

% Handling <> sytax for font sizes
% ================================
% The following code is based on LaTeX (\DeclareFontShape). Only the
% code for <> is preserved (no functions), and a default value can be
% set before the first <>. In addition, here single values take
% precedende over ranges. The original work is by the LaTeX Team,
% under license LPPL.

\def\enit@ifnot@nil#1{%
  \def\enit@a{#1}%
  \ifx\enit@a\@nnil
    \expandafter\@gobble
  \else
    \expandafter\@firstofone
  \fi}

\def\enit@remove@to@nnil#1\@nnil{}
\def\enit@remove@angles#1>{\enit@simple@size}

\def\enit@simple@size#1<{%
  \if<#1<%
    \expandafter\enit@remove@angles
  \else
    \def\enit@c{#1}%
    \expandafter\enit@remove@to@nnil
  \fi}

\def\enit@extractrange#1<#2>{%
  \ifx\enit@c\relax
    \def\enit@c{#1}%
  \fi
  \enit@isrange#2->\@nil#2>}

\def\enit@isrange#1-#2\@nil{%
   \if>#2%
     \expandafter\enit@check@single
   \else
     \expandafter\enit@check@range
   \fi}

\def\enit@check@range#1-#2>#3<#4\@nnil{%
  \enit@ifnot@nil{#3}{%
    \def\enit@b{\enit@extractrange<#4\@nnil}%
    \upper@bound=%
      \enit@ifunset{enit@sizefor@#2}{0#2\p@}{\@nameuse{enit@sizefor@#2}\p@}%
          %%% usar count@
    \ifdim\upper@bound=\z@ \upper@bound\maxdimen \fi
    \ifdim\f@size\p@<\upper@bound
      \lower@bound=%
      \enit@ifunset{enit@sizefor@#1}{0#1\p@}{\@nameuse{enit@sizefor@#1}\p@}%
      \ifdim\f@size\p@<\lower@bound
      \else
         \enit@simple@size#3<#4\@nnil
      \fi
    \fi
    \enit@b}}

\def\enit@check@single#1>#2<#3\@nnil{%
  \def\enit@b{\enit@extractrange<#3\@nnil}%
  \ifdim\f@size\p@=
     \enit@ifunset{enit@sizefor@#1}{0#1\p@}{\@nameuse{enit@sizefor@#1}\p@}%
     \enit@simple@size#2<#3\@nnil
     \let\enit@d\enit@c
  \fi
  \enit@b}

\def\enit@try@size@range#1{%
  \def\enit@a{#1}%
  \let\enit@c\relax  % last in range
  \let\enit@d\relax  % single
  \expandafter\enit@extractrange\enit@a <-*>\@nil<\@nnil
  \ifx\enit@d\relax\else\let\enit@c\enit@d\fi}

% \enit@setlength is defined in the options section

% This ends the code adapted from latex (under LPPL).

\def\SetEnumitemSize#1#2{%
  {\let\selectfont\relax
   #2%
   \expandafter\xdef\csname enit@sizefor@#1\endcsname{\f@size}}}

\SetEnumitemSize{script}\scriptsize
\SetEnumitemSize{tiny}\tiny
\SetEnumitemSize{footnote}\footnotesize
\SetEnumitemSize{small}\small
\SetEnumitemSize{normal}\normalsize
\SetEnumitemSize{large}\large
\SetEnumitemSize{Large}\Large
\SetEnumitemSize{LARGE}\LARGE
\SetEnumitemSize{huge}\huge
\SetEnumitemSize{Huge}\Huge

% +=============================+
% |         LIST TYPES          |
% +=============================+
%
% Displayed lists
% ===============
% #1 #2 implicit

\def\enit@dylist{%
  \enit@align@right
  \list}

\def\enit@endlist{%
  \enit@after
  \endlist
  \ifx\enit@series\relax\else % discards resume*, too
    \ifnum\enit@resuming=\@ne % ie, series=
      \enit@setresumekeys{series@\enit@series}\global\global
    \else % ie, resume=, resume*= (save count, but not keys)
      \enit@setresumekeys{series@\enit@series}\@gobblefour\global
    \fi
    \enit@afterlist
  \fi
  \ifnum\enit@resuming=\thr@@ % ie, resume* list (save count only)
    \enit@setresumekeys\@currenvir\@gobblefour\global
  \else
    \enit@setresumekeys\@currenvir\@empty\@empty
  \fi
  \aftergroup\enit@afterlist}

% #1 = either \@currenvir or series@<series>
% #2(keys) #3(counter) are \global, \@gobblefour or \@empty

\def\enit@setresumekeys#1#2#3{%
  \enit@toks\expandafter{\enit@savekeys}%
  \xdef\enit@afterlist{%
    #2\def\enit@noexcs{enit@resumekeys@#1}{\the\enit@toks}%
    \ifnum\enit@type=\z@ % ie, enum
      #3\def\enit@noexcs{enit@resume@#1}{%
        \csname c@\@listctr\endcsname
        \the\csname c@\@listctr\endcsname}%
    \fi}}

% Inline lists
% ============

% Definition of \@trivlist inside inline lists.  So, when
% \@trivlist is found in any displayed list (including quote,
% center, verbatim...) the default \@item is restored.

\def\enit@intrivlist{%
  \enit@changed@itemfalse
  \let\@item\enit@outer@item
  \let\par\@@par
  \let\@trivlist\enit@outer@triv
  \@trivlist}

% Keep track of \@item and \item changes

\newif\ifenit@changed@item
\enit@changed@itemfalse

\newif\ifenit@changeditem
\enit@changeditemfalse

% List
% ----

% Arguments, as before:
% \enitdp@<name>, <name>, <max-depth>, <format>
% About @newlist, see @initem.

\def\enit@inlist#1#2{%
  \ifnum\@listdepth>\enit@listdepth\relax
    \@toodeep
  \else
    \global\advance\@listdepth\@ne
  \fi
  \let\enit@align\@firstofone
  \def\@itemlabel{#1}%
  \@nmbrlistfalse
  \ifenit@changed@item\else
    \enit@changed@itemtrue
    \let\enit@outer@triv\@trivlist
    \let\@trivlist\enit@intrivlist
    \@setpar\@empty
    \let\enit@outer@item\@item
  \fi
  #2\relax
  \global\@newlisttrue
  \ifenit@boxmode
    \ifenit@changeditem\else
      \enit@changeditemtrue
      \let\enit@outeritem\item
    \fi
    \let\@item\enit@boxitem
  \else
    \let\@item\enit@noboxitem
    \ifx\enit@itemjoin@s\relax\else
      \PackageWarning{enumitem}%
         {itemjoin* discarded in mode unboxed\MessageBreak}%
    \fi
  \fi
  \let\enit@calcleft\relax
  \let\enit@afteritem\relax
  \ifenit@boxmode
    \global\setbox\enit@inbox\hbox\bgroup\color@begingroup
      \let\item\enit@endinbox
  \fi
  \ignorespaces}

\def\enit@endinlist{%
  \ifenit@boxmode
      \unskip
      \xdef\enit@afteritem{%
        \ifhmode\spacefactor\the\spacefactor\relax\fi}%
      \color@endgroup
    \egroup
    \ifdim\wd\enit@inbox=\z@
      \enit@itemerror
    \else
      \ifenit@noinitem\else
        \ifhmode\unskip\fi
        \enit@ifunset{enit@itemjoin@s}%
          {\enit@itemjoin}%
          {\enit@itemjoin@s}%
      \fi
      \unhbox\@labels
      \enit@afterlabel
      \unhbox\enit@inbox
      \enit@afteritem
    \fi
  \else
    \unskip
    \if@newlist
      \enit@itemerror
    \fi
  \fi
  \enit@after
  \global\advance\@listdepth\m@ne
  \global\@inlabelfalse
  \if@newlist
    \global\@newlistfalse
    \@noitemerr
  \fi
  \ifx\enit@series\relax\else % discards resume* list, too
    \ifnum\enit@resuming=\@ne % ie, series
      \enit@setresumekeys{series@\enit@series}\global\global
    \else % ie, resume, resume* (save count, but not keys)
      \enit@setresumekeys{series@\enit@series}\@gobblefour\global
    \fi
    \enit@afterlist
  \fi
  \ifnum\enit@resuming=\thr@@ % ie, resume* list (save count only)
    \enit@setresumekeys\@currenvir\@gobblefour\global
  \else
    \enit@setresumekeys\@currenvir\@empty\@empty
  \fi
  \aftergroup\enit@afterlist}

% \@item: unboxed
% ---------------

\def\enit@noboxitem[#1]{%
  \if@newlist
    \leavevmode % ships pending labels out
    \global\@newlistfalse
  \else
    \ifhmode
      \unskip
      \enit@itemjoin
    \else
      \noindent
    \fi
  \fi
  \if@noitemarg
    \@noitemargfalse
    \if@nmbrlist
      \refstepcounter{\@listctr}% after \unskip (hyperref)
    \fi
  \fi
  \mbox{\makelabel{#1}}%
  \enit@afterlabel
  \ignorespaces}

% \@item: boxed
% ------------
%
% We don't want \item to be executed locally, because it sets a flag
% (and hyperref adds another flag, too).  So, we redefine it inside
% the box to \enit@endinbox which ends the box and then use the actual
% (outer) \item.  labels are stored in another box, to detect empty
% boxes, ie, misplaced \item's.  Note the 2nd \item ends collecting
% the 1st item and ships it out, while the 3rd \item ends collecting
% the 2nd item, puts the itemjoin and then ships the 2nd item out.
% The flag enit@noinitem keeps track of that.

\newif\ifenit@noinitem

\def\enit@endinbox{%
    \unskip
    \xdef\enit@afteritem{%
      \ifhmode\spacefactor\the\spacefactor\relax\fi}%
    \color@endgroup
  \egroup
  \enit@outeritem}

\def\enit@boxitem[#1]{%
  \if@newlist
    \global\@newlistfalse
    \ifdim\wd\enit@inbox>\z@
       \enit@itemerror
    \fi
    \enit@noinitemtrue
    \leavevmode % ships pending labels out
  \else
    \ifdim\wd\enit@inbox=\z@
      \enit@itemerror
    \else
      \ifenit@noinitem
        \enit@noinitemfalse
      \else
        \ifhmode\unskip\fi
        \enit@itemjoin
      \fi
      \unhbox\@labels
      \enit@afterlabel
      \unhbox\enit@inbox
      \enit@afteritem
    \fi
  \fi
  \if@noitemarg
    \@noitemargfalse
    \if@nmbrlist
      \refstepcounter{\@listctr}%
    \fi
  \fi
  \sbox\@labels{\makelabel{#1}}%
  \let\enit@afteritem\relax
  \setbox\enit@inbox\hbox\bgroup\color@begingroup
    \let\item\enit@endinbox
    \hskip1sp % in case the first thing is \label
    \ignorespaces}

% Pause item
% ----------
%
% To do.
%
% The three types
% ===============
%
% enumerate and enumerate*
% ------------------------
%
% The following has 4 arguments, which in enumerate are:
% \@enumdepth, enum, \thr@@, <format>.
% In user defined environments they are:
% \enitdp@<name>, <name>, <max-depth>, <format>

\def\enit@enumerate{%
  \let\enit@list\enit@dylist
  \enit@enumerate@i}

\@namedef{enit@enumerate*}{%
  \let\enit@list\enit@inlist
  \enit@enumerate@i}

\def\enit@enumerate@i#1#2#3#4{%
  \ifnum#1>#3\relax
    \enit@toodeep
  \else
    \enit@prelist\z@{#1}{#2}%
    \edef\@enumctr{#2\romannumeral#1}%
    \expandafter
    \enit@list
      \csname label\@enumctr\endcsname
      {\usecounter\@enumctr
       \let\enit@calc\z@
       \def\makelabel##1{\enit@align{\enit@format{##1}}}%
       \enit@preset{#2}{#1}{#4}%
       \enit@normlabel\@itemlabel\@itemlabel
       \enit@ref
       \enit@calcleft
       \enit@before
       \enit@negwidth}%
    \enit@keyfirst
  \fi}

\let\enit@endenumerate\enit@endlist
\@namedef{enit@endenumerate*}{\enit@endinlist}

% itemize and itemize*
% --------------------
%
% The following has 4 arguments, which in itemize are:
% \@itemdepth, item, \thr@@, <format>.
% In user defined environments they are:
% \enitdp@<name>, <name>, <max-depth>, <format>

\def\enit@itemize{%
  \let\enit@list\enit@dylist
  \enit@itemize@i}

\@namedef{enit@itemize*}{%
  \let\enit@list\enit@inlist
  \enit@itemize@i}

\def\enit@itemize@i#1#2#3#4{%
  \ifnum#1>#3\relax
    \enit@toodeep
  \else
    \enit@prelist\@ne{#1}{#2}%
    \edef\@itemitem{label#2\romannumeral#1}%
    \expandafter
    \enit@list
      \csname\@itemitem\endcsname
       {\let\enit@calc\z@
        \def\makelabel##1{\enit@align{\enit@format{##1}}}%
        \enit@preset{#2}{#1}{#4}%
        \enit@calcleft
        \enit@before
        \enit@negwidth}%
    \enit@keyfirst
  \fi}

\let\enit@enditemize\enit@endlist
\@namedef{enit@enditemize*}{\enit@endinlist}

% description and description*
% ----------------------------
%
% Make sure \descriptionlabel exists:

\providecommand*\descriptionlabel[1]{%
  \hspace\labelsep
  \normalfont\bfseries#1}

\@namedef{enit@description*}{%
  \let\enit@list\enit@inlist
  \enit@description@i}

\def\enit@description{%
  \let\enit@list\enit@dylist
  \enit@description@i}

\def\enit@description@i#1#2#3#4{%
  \ifnum#1>#3\relax
    \enit@toodeep
  \else
    \enit@list{}%
      {\let\enit@type\tw@
       \advance#1\@ne
       \labelwidth\z@
       \enit@align@left
       \let\makelabel\descriptionlabel
       \enit@style@standard
       \enit@preset{#2}{#1}{#4}%
       \enit@calcleft
       \let\enit@svlabel\makelabel
       \def\makelabel##1{%
         \labelsep\z@
         \ifenit@boxdesc
           \enit@svlabel{\enit@align{\enit@format{##1}}}%
         \else
           \nobreak
           \enit@svlabel{\enit@format{##1}}%
           \aftergroup\enit@postlabel
         \fi}%
       \enit@before
       \enit@negwidth}%
     \enit@keyfirst
  \fi}

\let\enit@enddescription\enit@endlist
\@namedef{enit@enddescription*}{\enit@endinlist}

% trivlist
% ========

\def\enit@trivlist{%
  \let\enit@type\tw@
  \parsep\parskip
  \csname @list\romannumeral\the\@listdepth\endcsname
  \@nmbrlistfalse
  \enit@tagid
  \enit@setglobalkeys % ie, list and list<num>
  \enit@setkeys{trivlist}%
  \enit@setkeys{trivlist\romannumeral\@listdepth}%
  \@trivlist
  \labelwidth\z@
  \leftmargin\z@
  \itemindent\z@
  \let\@itemlabel\@empty
  \def\makelabel##1{##1}}

% Description styles
% ==================
%
% the next definition is somewhat tricky because labels are boxed.
% That's fine when the label is just placed at the begining of a line
% of text, but when the box is placed without horizontal material,
% leading is killed.  So, we need change somehow \box to \unhbox, but
% I don't want to modify \@item.  The code below presumes \@item has
% not been changed and arguments gobble the part setting \@labels,
% which is replaced by a new one.
%
% The default value in description is itemindent=!, but some styles
% (those whose item text begin at a fixed place, ie, nextline,
% multiline and sameline) change it to labelwidth=!.
%
% We must be careful with the group and the whatsit added by color to
% boxes.  Alignment is applied here and some adjustments in skips are
% necessary to get proper line breaks (including a \nobreak at the
% beginning of \enit@align, ie, after the first whatsit, see above).
% To "pass" the inner group added by color to the box, \enit@postlabel
% ckecks if the following is }.  ie, \egroup -- if not, the box has
% not reached yet its end.

\def\enit@postlabel{%
  \@ifnextchar\egroup
    {\aftergroup\enit@postlabel}%
    {\enit@postlabel@i}}

\def\enit@postlabel@i#1#2#3#4#5{%
  \def\enit@lblpenalty{\penalty\z@\hskip\skip@}%
  \ifenit@nextline
    \ifdim\wd\@tempboxa>\labelwidth
      \def\enit@lblpenalty{\newline\@nobreaktrue}%
    \fi
  \fi
  \everypar{%
    \@minipagefalse
    \global\@newlistfalse
    \if@inlabel
      \global\@inlabelfalse
      {\setbox\z@\lastbox
       \ifvoid\z@
         \kern-\itemindent
       \fi}%
      \unhbox\@labels
      \skip@\lastskip % Save last \labelsep
      \unskip % Remove it
      \enit@lblpenalty % Restore it, after penalty
    \fi
    \if@nobreak
      \@nobreakfalse
      \clubpenalty\@M
    \else
      \clubpenalty\@clubpenalty
      \everypar{}%
    \fi}%
  \def\enit@a{#1#2#3#4}%
  \def\enit@b{\global\setbox\@labels\hbox}%
  \ifx\enit@a\enit@b\else
    \enit@error{Non standard \string\item}%
      {A class or a package has redefined \string\item\MessageBreak
       and I do not know how to continue}%
  \fi
  \global\setbox\@labels\hbox{%
    \unhbox\@labels
    \hskip\itemindent
    \hskip-\labelwidth
    \hskip-\labelsep
    \ifdim\wd\@tempboxa>\labelwidth
      \enit@align{\unhbox\@tempboxa}\unskip % Removes (typically) \hfil
    \else
      \leavevmode\hbox to\labelwidth{\enit@align{\unhbox\@tempboxa}}%
    \fi
    \hskip\labelsep}}

% +=============================+
% |     (RE)DEFINING LISTS      |
% +=============================+
%
% Set keys/values
% ===============
% Remember \romannumeral0 expands to nothing.
% #1 = list name, #2 = level, #3 = flag if star, #4 = keys/values

\let\enit@keys@sizes\relax

\def\enit@saveset#1#2#3#4{%
  \setcounter{enit@cnt}{#2}%
  \ifx\enit@forsize\@empty
    \ifcase#3%
      \expandafter
      \def\csname enit@@#1\romannumeral\c@enit@cnt\endcsname{#4}%
    \or
      \expandafter\let\expandafter\enit@b
        \csname enit@@#1\romannumeral\c@enit@cnt\endcsname
      \ifx\enit@b\relax
        \let\enit@b\@empty
      \fi
      \expandafter\def
        \csname enit@@#1\romannumeral\c@enit@cnt\expandafter\endcsname
        \expandafter{\enit@b,#4}%
    \fi
  \else
    \ifcase#3%
      \enit@ifunset{enit@@#1\romannumeral\c@enit@cnt}%
        {\expandafter\let
         \csname enit@@#1\romannumeral\c@enit@cnt\endcsname\@empty}%
        {}%
      \expandafter\let\expandafter\enit@b
        \csname enit@@#1\romannumeral\c@enit@cnt @@sizes\endcsname
      \ifx\enit@b\relax
        \let\enit@b\@empty
      \fi
      \toks@\expandafter{\enit@b}%
      \edef\enit@b{\the\toks@\enit@forsize\enit@keys@sizes}%
      \expandafter\def
        \csname enit@@#1\romannumeral\c@enit@cnt @@sizes\expandafter\endcsname
        \expandafter{\enit@b{#4}}%
    \else
      \enit@error{* and \string<\string> are not compatible}%
        {Use either * or angles, but not both.}%
    \fi
  \fi}

% TODO: more robust tests (catch wrong key names, but not easy)

% Internally, LaTeX uses a short name for enumerate (enum)
% and itemize (item). To be consistent with this convention,
% a couple of macros provide a "translation". I'm not very
% happy with the current implementation.

\def\enit@shortenumerate{enum}
\def\enit@shortitemize{item}

\newcommand\setlist{%
  \@ifstar{\enit@setlist\@ne}{\enit@setlist\z@}}

\def\enit@setlist#1{%
  \@ifnextchar<%
    {\enit@setlist@q#1}%
    {\let\enit@forsize\@empty\enit@setlist@n#1}}

% Default definitions. Modified below with option 'sizes':

\def\enit@setlist@q#1<#2>{%
  \enit@error
    {Activate this feature with options 'sizes'}%
    {Size dependent setting with \string<\string> must be\MessageBreak
     explicitly activated with the package option 'sizes'}}

\def\enit@setlist@n#1{%
  \@ifnextchar[{\enit@setlist@x#1}{\enit@setlist@i#1\@empty}}

% Let's accept \setlist[]*{}, too, because an error in <=3.5.1

\def\enit@setlist@x#1[#2]{%
  \@ifstar{\enit@setlist@i\@ne{#2}}{\enit@setlist@i#1{#2}}}

% #1 list names/levels, #2 keys/values

% #1 star flag, #2 list names/levels, #3 keys/values

\def\enit@setlist@i#1#2#3{%
  \let\enit@eltnames\relax
  \let\enit@b\@empty
  \let\enit@eltlevels\relax
  \let\enit@c\@empty
  \protected@edef\enit@a{#2}%
  \@for\enit@a:=\enit@a\do{% the 2nd enit@a is first expanded
    \enit@ifunset{enitdp@\enit@meaning\enit@a}%
      {\edef\enit@c{\enit@c\enit@eltlevels{\enit@a}}}%
      {\enit@ifunset{enit@short\enit@meaning\enit@a}%
         \@empty
         {\edef\enit@a{\@nameuse{enit@short\enit@a}}}%
       \edef\enit@b{\enit@b\enit@eltnames{\enit@a}}}}%
  \ifx\enit@b\@empty
     \def\enit@b{\enit@eltnames{list}}%
  \fi
  \ifx\enit@c\@empty
     \def\enit@c{\enit@eltlevels{0}}%
  \fi
  \def\enit@eltnames##1{%
    \def\enit@a{##1}%
    \enit@c}%
  \def\enit@eltlevels##1{%
    \enit@saveset\enit@a{##1}#1{#3}}%
  \enit@b}%

% Deprecated:

\newcommand\setdisplayed[1][0]{\setlist[trivlist,#1]}
\let\enitdp@trivlist\@empty % dummy, let know it exists
\newcommand\setenumerate[1][0]{\setlist[enumerate,#1]}
\newcommand\setitemize[1][0]{\setlist[itemize,#1]}
\newcommand\setdescription[1][0]{\setlist[description,#1]}

% New lists
% =========

% When defining a list, \label... and counters must be defined
% for each level, too:

\def\enit@xset@itemize{%
  \@namedef{label\enit@c\romannumeral\count@}{%
    \enit@error{Undefined label}%
      {You have defined a list, but labels have
       not been setup.\MessageBreak
       You can set the label field with \string\setlist.}}}
\@namedef{enit@xset@itemize*}{\enit@xset@itemize}

\def\enit@xset@enumerate{%
  \enit@xset@itemize
  \enit@ifunset{c@\enit@c\romannumeral\count@}%
    {\@definecounter{\enit@c\romannumeral\count@}}{}}
\@namedef{enit@xset@enumerate*}{\enit@xset@enumerate}

\let\enit@xset@description\@empty
\@namedef{enit@xset@description*}{\enit@xset@description}

\newcommand\newlist{\enit@newlist\newenvironment}
\newcommand\renewlist{\enit@newlist\renewenvironment}

% <new/renew>, <name>, <type>, <max-depth>

\def\enit@newlist#1#2#3#4{%
  \enit@ifunset{enit@xset@#3}%
    {\enit@error{Unknown list type '#3')}%
          {Valid types are:
           enumerate, itemize, description,\MessageBreak
           enumerate*, itemize*, description*}}%
    {}%
  \setcounter{enit@cnt}{#4}%
  \count@\@ne
  \enit@ifunset{enit@short#2}%
    {\def\enit@c{#2}}%
    {\edef\enit@c{\csname enit@short#2\endcsname}}%
  \loop
    \@nameuse{enit@xset@#3}% Uses \enit@c
    \ifnum\count@<\c@enit@cnt
    \advance\count@\@ne
  \repeat
  \enit@ifunset{enitdp@#2}%
    {\expandafter\newcount\csname enitdp@#2\endcsname}{}%
  \csname enitdp@#2\endcsname\z@
  \advance\c@enit@cnt\m@ne
  \edef\enit@a{%
    \noexpand#1{#2}[1][]%
      {\enit@noexcs{enit@#3}%
       \enit@noexcs{enitdp@#2}%
       {\enit@c}%
       {\the\c@enit@cnt}%
       {####1}}%
      {\enit@noexcs{enit@end#3}}}%
  \enit@a}

% Changing the default nesting limit
% ----------------------------------

\newcommand\setlistdepth{\def\enit@listdepth}
\setlistdepth{5}

% +=============================+
% |       PACKAGE OPTIONS       |
% +=============================+

\newif\ifenit@loadonly

\DeclareOption{ignoredisplayed}{\let\enit@trivlist\trivlist}
\DeclareOption{includedisplayed}{%
  \def\enit@setglobalkeys{%
    \enit@setkeys{list}%
    \enit@setkeys{list\romannumeral\@listdepth}}}
\let\enit@setglobalkeys\relax

\DeclareOption{loadonly}{\enit@loadonlytrue}

\DeclareOption{shortlabels}
  {\def\enit@shl#1{%
     \ifnum\enit@type=\tw@
       \enit@toks{#1}%
     \else
       \def\enit@c{#1}%
       \enit@first#1,\@nil\@@nil % Returns enit@toks
    \fi}}

\DeclareOption{inline}
  {\newenvironment{enumerate*}[1][]%
     {\@nameuse{enit@enumerate*}\enitdp@enumerate{enum}\thr@@{#1}}
     {\@nameuse{enit@endenumerate*}}
   \newenvironment{itemize*}[1][]%
     {\@nameuse{enit@itemize*}\enitdp@itemize{item}\thr@@{#1}}
     {\@nameuse{enit@enditemize*}}
   \newenvironment{description*}[1][]%
     {\@nameuse{enit@description*}\enitdp@description{description}\@M{#1}}
     {\@nameuse{enit@enddescription*}}}

\def\enit@setlength{\setlength}
\DeclareOption{sizes}
  {\def\enit@setlength#1#2{%
     \enit@try@size@range{#2}%  Returns \enit@c
     \setlength#1{\enit@c}}%
   \def\enit@setkeys#1{%
     \enit@ifunset{enit@@#1}{}%
       {\let\enit@c\@empty
        \enit@ifunset{enit@@#1@@sizes}{}%
          {\expandafter\let\expandafter\enit@a\csname enit@@#1@@sizes\endcsname
           \expandafter\enit@try@size@range\expandafter{\enit@a}%
           \def\enit@keys@sizes{\def\enit@c}%
           \enit@c
           \let\enit@keys@sizes\relax}%
        \expandafter\expandafter\expandafter
          \enit@setkeys@i
          \csname enit@@#1\expandafter\endcsname
          \expandafter,\enit@c\@@}}%
   \def\enit@setlist@q#1<#2>{%
     \def\enit@forsize{<#2>}%
     \enit@setlist@n#1}}

\chardef\enit@seriesopt\z@
\DeclareOption{series=override}{\chardef\enit@seriesopt\tw@}

\let\enit@shl\enit@toks

\ProcessOptions

\let\trivlist\enit@trivlist

% If there is no loadonly, redefine the basic lists:

\ifenit@loadonly\else

\let\enitdp@enumerate\@enumdepth
\renewenvironment{enumerate}[1][]
  {\enit@enumerate\enitdp@enumerate{enum}\thr@@{#1}}
  {\enit@endenumerate}

\let\enitdp@itemize\@itemdepth
\renewenvironment{itemize}[1][]
  {\enit@itemize\enitdp@itemize{item}\thr@@{#1}}
  {\enit@enditemize}

\newcount\enitdp@description
\renewenvironment{description}[1][]
  {\enit@description\enitdp@description{description}\@M{#1}}
  {\enit@enddescription}

\fi

% +=============================+
% |            TOOLS            |
% +=============================+

\def\enit@drawrule#1#2#3#4{%
  \rlap{%
    \ifdim#1>0pt\relax
      \vrule width #1 height #2 depth -#3\relax
    \else\ifdim#1=0pt\relax
      %
    \else
      \hskip#1%
      \vrule width -#1 height #2 depth -#4\relax
    \fi\fi}}
  
\def\DrawEnumitemLabel{%
  \begingroup
    \item[]%
    \hskip-\labelsep
    \enit@drawrule\labelsep{4pt}{3pt}{2.3pt}%
    \hskip-\labelwidth
    \enit@drawrule\labelwidth{6pt}{5pt}{4.3pt}%
    \hskip\labelwidth
    \hskip\labelsep
    %
    \hskip-\itemindent
    \enit@drawrule\itemindent{2pt}{1pt}{.3pt}%
    \rlap{\vrule height 9pt depth .5pt}%
    \hskip-\leftmargin
    \rlap{\vrule height 9pt depth .5pt}%
    \enit@drawrule\labelindent{8pt}{7pt}{6.5pt}%
    % \message{\the\labelindent/\the\labelwidth/\the\labelsep/\the\itemindent}%
  \endgroup} 

% TODO -  option 'verbose'

% +=============================+
% |        TWO-PASS TOOLS       |
% +=============================+

% TODO - Not for the moment, because both tools would require to take
% into account series and resume, which is not simple. Also, are they
% applied to all lists, by type, by level, by name, etc.? Document how
% to do it in at least the simplest cases.
%
% - reverse counting
% - automatic widest, based on existing labels.

\endinput

MIT License
-----------

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
`
)
