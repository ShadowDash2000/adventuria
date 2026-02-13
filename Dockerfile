FROM golang:1.26.0-alpine AS builder

RUN mkdir -p "/adventuria/"
COPY . /adventuria/

WORKDIR /adventuria/
RUN go mod download

RUN apk update \
    && apk upgrade

RUN go build -a -installsuffix cgo -o ./adventuria cmd/main.go

EXPOSE 8080

CMD ["./adventuria", "serve", "--http=0.0.0.0:8080"]