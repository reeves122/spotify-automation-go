FROM golang:1.17 AS builder

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN go build -o spotify-automation-go main.go


FROM debian:buster-slim

RUN apt-get update && apt-get install -y --no-install-recommends apt-utils ca-certificates

COPY --from=builder /app/spotify-automation-go /

CMD ["/spotify-automation-go"]