FROM golang:latest
WORKDIR /api

COPY . .
RUN go mod download
RUN go build

ENV PORT=8081

EXPOSE 8080

ENTRYPOINT ./LetsTalk