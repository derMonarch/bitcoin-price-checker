FROM golang:1.23.4

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o main ./cmd/btcprice.go

EXPOSE 8080

CMD ["/app/main"]