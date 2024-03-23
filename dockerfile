FROM golang:latest

WORKDIR /GoApi

COPY ./ ./
COPY go.mod go.sum ./
RUN go mod download


RUN go build -o goapi ./main.go

CMD ./goapi