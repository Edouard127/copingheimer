FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./src/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]