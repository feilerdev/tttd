# Container image that runs your code
FROM golang:1.26-bookworm

WORKDIR /go/src/action
COPY . .

RUN go mod tidy
RUN go build -o /go/bin/action

ENTRYPOINT ["/go/bin/action"]
