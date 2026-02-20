FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY telegram.go .env ./

RUN ["go","run","telegram.go"]
