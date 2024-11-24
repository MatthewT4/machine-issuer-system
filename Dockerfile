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

ENV DB_URL postgresql://postgres:postgres@postgres:5432/postgres
ENV AUTH_SECRET_KEY lskjfijojslkfjsoidjfs
ENV AUTH_TTL 24
ENV AUTH_COOKIE_NAME session_token
ENV SSH_FILE_PATH /Users/prc-94/.ssh/hack_ed

CMD ["./main"]
