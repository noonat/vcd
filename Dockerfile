# This Dockerfile is intended to be used for development.

FROM golang:1.10

RUN go get -u github.com/githubnemo/CompileDaemon

COPY . /go/src/github.com/noonat/vcd
WORKDIR /go/src/github.com/noonat/vcd
RUN go install github.com/noonat/vcd/cmd/vcd

EXPOSE 8080
CMD ["vcd"]
