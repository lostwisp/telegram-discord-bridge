FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY telegram.go .env ./

CMD ["go","run","telegram.go"]
