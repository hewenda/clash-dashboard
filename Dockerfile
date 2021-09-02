ARG GO_VERSION=1.16.7

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update \
    && apk add --no-cache build-base git
RUN git clone https://github.com/hewenda/clash-dashboard.git /app

ENV GIN_MODE="release"
ENV RUN_ENV="production"

WORKDIR /app

RUN go mod download
RUN GOARCH=wasm GOOS=js go build -o web/app.wasm
RUN go build .

EXPOSE 3000

ENTRYPOINT ["/app/clash"]