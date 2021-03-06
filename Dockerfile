FROM golang:latest as builder

WORKDIR /build
ADD . /build

RUN CGO_ENABLED=0 go build -o google-storage-proxy ./cmd/

FROM alpine:latest
LABEL org.opencontainers.image.source=https://github.com/cirruslabs/google-storage-proxy/

WORKDIR /svc
COPY --from=builder /build/google-storage-proxy /svc/
ENTRYPOINT ["/svc/google-storage-proxy"]
