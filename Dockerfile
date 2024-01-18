FROM golang:1.21-alpine

COPY . /home/src
WORKDIR /home/src
RUN go build - /bin/action ./cmd/main.go

ENTRYPOINT [ "/bin/action" ]
