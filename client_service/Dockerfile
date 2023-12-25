FROM golang:alpine
WORKDIR /build
COPY . .
RUN go build -o client-service ./cmd/main.go
CMD ["./bank-app"]
