FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o socket ./src
RUN chmod +x socket
EXPOSE 29229
EXPOSE 80
CMD ["./socket"]