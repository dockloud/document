FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/document ./cmd/app
RUN CGO_ENABLED=0 go build -o /app/migrate ./cmd/migrate

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/document ./document
COPY --from=builder /app/migrate ./migrate

EXPOSE 50160

CMD ["./document"]