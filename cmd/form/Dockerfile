FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o form -mod=vendor -ldflags='-s -w'  -installsuffix cgo ./cmd/form/main.go

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /form
COPY --from=builder ./build/form ./cmd/
