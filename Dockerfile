FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o wow-tracker ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/wow-tracker .
EXPOSE 8080
CMD ["./wow-tracker"]
