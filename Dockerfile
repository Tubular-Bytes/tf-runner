FROM registry.0x42.in/terrence/golang:1.24-alpine3.21 AS builder

RUN apk add git make

WORKDIR /build
COPY . .
RUN make build

FROM registry.0x42.in/terrence/base:alpine3.21
COPY --from=builder /build/bin/runner /usr/bin/runner

ENTRYPOINT ["/usr/bin/runner"]