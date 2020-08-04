FROM golang
WORKDIR /app
RUN go build
ENTRYPOINT [ "./letstalk" ]