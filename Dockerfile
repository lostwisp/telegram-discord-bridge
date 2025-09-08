FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go get -u ./... &&  go get github.com/bwmarrin/discordgo

COPY . .

ENTRYPOINT ["go", "run", "cmd/main.go"]
