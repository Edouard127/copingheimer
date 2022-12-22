FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o worker .
RUN chmod +x worker
CMD ["./worker"]