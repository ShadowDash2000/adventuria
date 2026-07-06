FROM golang:1.26.2-alpine AS builder

WORKDIR /adventuria

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o adventuria cmd/main.go

FROM alpine:latest

WORKDIR /adventuria

COPY --from=builder /adventuria/adventuria .

EXPOSE 8080

CMD ["./adventuria", "serve", "--http=0.0.0.0:8080"]