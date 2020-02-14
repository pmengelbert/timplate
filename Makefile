build:
	curl http://ctan.math.washington.edu/tex-archive/macros/latex/contrib/enumitem/enumitem.sty -s > enumitem.sty
	/usr/bin/go build .

install:
	cp timplate /usr/local/bin
	mkdir -p sudo mkdir /usr/share/texmf/tex/latex/enumitem || true
	cp enumitem.sty /usr/share/texmf/tex/latex/enumitem/
