FROM alpine as certs
RUN apk update && apk add ca-certificates

FROM golang:1.16.6-alpine3.14 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o job -mod=vendor -ldflags='-s -w'  -installsuffix cgo ./pkg/job/v0.0.1/.

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs

WORKDIR /job
COPY --from=builder ./build/job ./cmd/


ENTRYPOINT ["./cmd/job","-config=/configs/config.yml"]