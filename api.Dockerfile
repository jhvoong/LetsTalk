FROM golang:latest
WORKDIR /api

COPY go.mod go.sum main.go ./
RUN go mod download
RUN go build

EXPOSE 8080

ENTRYPOINT ./LetsTalk