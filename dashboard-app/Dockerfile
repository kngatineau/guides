FROM node:18-alpine3.18 as ui

RUN npm i -g pnpm

WORKDIR /src

COPY ui/package.json .
COPY ui/pnpm-lock.yaml .

RUN pnpm i

COPY ui/. .

RUN pnpm build

FROM golang:1.21-alpine3.18 as builder

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd cmd

RUN mkdir -p ui
COPY --from=ui /src/dist ui/dist
COPY ui/embed.go ui/embed.go

RUN mkdir -p bin

RUN go build -o bin/api ./cmd/api/...

FROM alpine:3.18

COPY --from=builder /src/bin/api /usr/bin/api

CMD /usr/bin/api
