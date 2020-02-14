FROM golang:1.13 AS build
WORKDIR /usr/local/bin
COPY . .
ENV CGO_ENABLED 0 
RUN go build .

FROM alpine:latest
COPY --from=build /usr/local/bin/timplate /app/timplate
RUN apk add curl make go texlive
ENV PATH="/app:${PATH}"
WORKDIR /lolz
RUN curl http://ctan.math.washington.edu/tex-archive/macros/latex/contrib/enumitem/enumitem.sty -s > enumitem.sty
RUN mkdir -p /root/texmf/tex/latex/mystuff
RUN mv enumitem.sty /root/texmf/tex/latex/mystuff/enumitem.sty
ENTRYPOINT ["/app/timplate"]
