FROM golang:1.17 as builder

WORKDIR /app/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY main.go .
COPY splunk.go .
RUN CGO_ENABLED=0 go build -o /main

ENTRYPOINT ["/main"]
