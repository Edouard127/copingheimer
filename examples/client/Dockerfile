FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server .
RUN chmod +x server
EXPOSE 29969
CMD ["./server"]