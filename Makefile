export PATH := /usr/bin:$(PATH)
export HOME := $(HOME)

build:
	mkdir bin || true
	go build -o bin/timplate . 

install:
	mv bin/timplate /usr/local/bin/timplate
