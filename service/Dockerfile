FROM golang:1.26

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG ENTRY_POINT
RUN go build -v -o /usr/local/bin/app "${ENTRY_POINT}"

CMD ["app"]