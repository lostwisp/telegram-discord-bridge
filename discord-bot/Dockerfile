FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY discord.go .env ./

CMD ["go","run","discord.go"]