FROM golang:bookworm AS builder

WORKDIR /build
COPY . .
RUN make build

FROM debian:bookworm
COPY --from=builder /build/bin/runner /usr/bin/runner

ENTRYPOINT ["/usr/bin/runner"]