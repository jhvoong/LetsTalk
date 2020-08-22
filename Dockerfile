FROM golang:latest
WORKDIR /app

COPY backend backend

COPY go.mod go.sum main.go ./
RUN go mod download

RUN go build
ENTRYPOINT ./LetsTalk
EXPOSE 8080