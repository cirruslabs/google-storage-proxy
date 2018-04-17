FROM golang:1.10-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /go/src/github.com/cirruslabs/google-storage-proxy
ADD . /go/src/github.com/cirruslabs/google-storage-proxy

RUN go get -t -v ./... && \
    go test -v ./... && \
    go build -o google-storage-proxy ./cmd/

FROM alpine
WORKDIR /svc
COPY --from=builder /go/src/github.com/cirruslabs/google-storage-proxy/google-storage-proxy /svc/
ENTRYPOINT ["/svc/google-storage-proxy"]