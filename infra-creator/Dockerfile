FROM golang:1.18

RUN mkdir /app
WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

WORKDIR /app

RUN go build

ENTRYPOINT ["./gitpod-crossplane-installer"]
