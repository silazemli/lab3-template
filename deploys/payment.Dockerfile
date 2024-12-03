FROM golang:1.23.1 AS builder

WORKDIR /app

COPY ["./go.mod", "./go.sum", "./"]
RUN go mod download

COPY .. ./

RUN go build -o /app/main ./internal/cmd/payment/main.go
RUN chmod +x /app/main

EXPOSE 8060
CMD /app/main