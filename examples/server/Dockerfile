FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o client .
RUN chmod +x server
CMD ["./client"]