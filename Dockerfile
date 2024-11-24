FROM golang:1.23.3-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

FROM golang:1.23.3 as runner

COPY --from=builder /app/main /app/main

WORKDIR /app

EXPOSE 80

CMD ["./main"]
