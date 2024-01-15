ARG GO_VERSION=1.21

FROM golang:1.21 AS builder

ARG GO111MODULE=on
ARG GOSUMDB=off

WORKDIR /go/src/all-wallet

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/all-wallet ./cmd/all-wallet


FROM alpine as app

WORKDIR /app

COPY --from=builder /go/src/all-wallet/bin/all-wallet .
COPY config.yaml config.yaml

ENTRYPOINT ["/app/all-wallet"]

