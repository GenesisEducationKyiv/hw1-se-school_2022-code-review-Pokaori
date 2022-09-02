FROM golang:1.13-alpine AS build_base

WORKDIR ~/bitcoin-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /bitcoin-service


FROM alpine:latest

COPY --from=build_base /bitcoin-service /app/bitcoin-service
COPY .env /.env

EXPOSE 8000

CMD ["/app/bitcoin-service"]