FROM golang:latest
WORKDIR /api

COPY . .
RUN go mod download
RUN go build

ENV PORT=8080

EXPOSE ${PORT}

ENTRYPOINT ./LetsTalk