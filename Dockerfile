# syntax=docker/dockerfile:1

FROM golang:1.20-alpine as Builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /app/subspace

FROM alpine

WORKDIR /app

COPY --from=Builder /app/sql /app/sql
COPY --from=Builder /app/subspace /app/subspace

ENTRYPOINT ["./subspace"]