FROM golang:1.21-alpine

COPY . /home/src
WORKDIR /home/src
RUN go build -o /bin/action ./cmd/action

ENTRYPOINT [ "/bin/action" ]
