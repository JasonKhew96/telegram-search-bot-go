# syntax=docker/dockerfile:1
FROM golang:1.22-bullseye AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o main .


FROM gcr.io/distroless/static-debian12
WORKDIR /app

COPY --from=builder /build/main /app/main

ENTRYPOINT [ "/app/main" ]
