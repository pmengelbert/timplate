export PATH := /usr/bin:$(PATH)
export HOME := $(HOME)

build:
	mkdir bin || true
	go build -o bin/timplate . 
	mkdir -p $(HOME)/texmf/tex/latex/enumitem
	mv enumitem.sty $(HOME)/texmf/tex/latex/enumitem

install:
	mv bin/timplate /usr/local/bin/timplate
	mkdir -p $(HOME)/texmf/tex/latex/enumitem
	mv enumitem.sty $(HOME)/texmf/tex/latex/enumitem
