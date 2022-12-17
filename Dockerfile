FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./src/
RUN chmod +x main
CMD ["./main"]