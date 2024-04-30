FROM golang:1.22.2

WORKDIR /usr/src/app

# RUN go install github.com/cosmtrek/air@latest

COPY . .

USER 10014


RUN go mod download

CMD ["go", "run", "cmd/main.go"]
