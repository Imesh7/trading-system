FROM golang:1.21.3

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .

USER 10014


RUN go mod tidy