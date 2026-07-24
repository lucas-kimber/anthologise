FROM golang:1.26 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG ENTRY_POINT
RUN CGO_ENABLED=0 go build -v -o /app "${ENTRY_POINT}"


FROM scratch

COPY --from=builder /app /app

ENTRYPOINT ["/app"]