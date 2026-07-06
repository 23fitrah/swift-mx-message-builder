FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o swift-mx-message-builder .

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/swift-mx-message-builder .
COPY --from=builder /app/.env .

RUN chmod +x ./swift-mx-message-builder

EXPOSE 8080

CMD ["./swift-mx-message-builder"]
