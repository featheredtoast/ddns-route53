FROM golang:1.16

WORKDIR /go/src/cli

COPY cli /go/src/cli
COPY internal /go/src/internal

RUN go get -d -v ./...
RUN go build -v -o /go/bin/ddns ./... 

ENTRYPOINT ["ddns", "continuous"]