FROM golang:1.23.5-alpine AS builder

WORKDIR /todo-list

COPY . .

RUN apk add --no-cache gcc g++ musl-dev libc6-compat

ENV CGO_ENABLED=1 CC="gcc"
RUN go build -o todo-list-app

FROM alpine:latest

WORKDIR /todo-list

RUN apk add --no-cache libc6-compat

COPY --from=builder /todo-list/ .

CMD ["./todo-list-app"]