ARG GO_VERSION=1.21
FROM golang:${GO_VERSION} AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN make dist

FROM scratch AS base-layer

STOPSIGNAL SIGINT

ENTRYPOINT ["/watchdog"]

FROM base-layer as watchdog

COPY --from=builder --chown=1000:1000 --chmod=755 /app/dist/watchdog /

USER 1000

FROM base-layer AS dev

COPY --chown=1000:1000 --chmod=755 ./dist/watchdog /

USER 1000