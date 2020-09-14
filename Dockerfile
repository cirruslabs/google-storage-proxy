FROM golang:latest as builder

WORKDIR /build
ADD . /build

RUN go get -t -v ./... && \
    go test -v ./... && \
    go build -o google-storage-proxy ./cmd/

FROM alpine:latest
LABEL org.opencontainers.image.source=https://github.com/cirruslabs/google-storage-proxy/

WORKDIR /svc
COPY --from=builder /build/google-storage-proxy /svc/
ENTRYPOINT ["/svc/google-storage-proxy"]