FROM golang:1.18 AS builder

RUN mkdir /app && mkdir /app/bin-out
ADD . /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /go-example

CMD ["/go-example"]
