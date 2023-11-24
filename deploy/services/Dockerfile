FROM golang:alpine AS base

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS test

RUN make lint

FROM base AS publish

RUN go build -o ./out/create ./cmd/create
RUN go build -o ./out/query ./cmd/query
RUN go build -o ./out/push ./cmd/push

FROM alpine:latest AS final

WORKDIR /usr/bin

COPY --from=publish ./src/out .
