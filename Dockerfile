FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /go/src/gateway

COPY . .

ENV GATEWAY_CONFIG_PATH=./config/config.yaml

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./gateway ./cmd/main.go

FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates

WORKDIR /root

ENV GATEWAY_CONFIG_PATH=./config/config.yaml

RUN mkdir -p /root/config

COPY --from=builder /go/src/gateway/config ./config

COPY --from=builder /go/src/gateway/gateway .

EXPOSE 80

ENTRYPOINT [ "/root/gateway" ]