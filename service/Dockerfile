# syntax=docker/dockerfile:1

FROM golang:1.26 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG ENTRY_POINT

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -v -o /app "${ENTRY_POINT}"


FROM scratch

COPY --from=builder /app /app

ENTRYPOINT ["/app"]