# syntax=docker/dockerfile:1
FROM golang:1.22-bullseye AS builder
WORKDIR /app
COPY . .

RUN apt-get update && \
	apt-get install -y sqlite3 && \
	rm -rf /var/lib/apt/lists/*

RUN go mod download
RUN go build -o main .
RUN ln -s config/config.yaml config.yaml && \
	ln -s config/data.db data.db
RUN rm -rf docker entity models *.go

ENTRYPOINT [ "/app/docker-entrypoint.sh" ]
