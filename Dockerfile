FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY src .
RUN go build -o socket .
RUN chmod +x socket
EXPOSE 29229
EXPOSE 80
CMD ["./socket"]