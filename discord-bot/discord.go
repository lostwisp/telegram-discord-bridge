package discord

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	tgdis "github.com/lostwisp/telegram-discord-bridge/gRPC/tg-dis"
	"google.golang.org/grpc"
)

type server struct {
	tgdis.UnimplementedTgdisMessageServiceServer
}

func (s *server) NewMessage(ctx context.Context, req *tgdis.MessageRequest) (*tgdis.MessageResponse, error) {
	return &tgdis.MessageResponse{Score: 0}, nil
}

var bot *discordgo.Session

func GetEnv() string {
	return os.Getenv("DISCORD_BOT_TOKEN")
}

func main() {
	//Создание сервера для прослушивания gRPC
	conn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Println(err)
	}

	grpcServer := grpc.NewServer()

	tgdis.RegisterTgdisMessageServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//Создание бота с WebSocket
	token := GetEnv()
	bot, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Println(err)
	}
	bot.AddHandler(pingPong)
	err = bot.Open()
	if err != nil {
		log.Println(err)
	}
	defer bot.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}

func pingPong(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Игнорируем сообщения от самого бота (чтобы избежать циклов)
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Если сообщение содержит "!ping", отвечаем "Pong!"
	if m.Content == "!ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong! 🏓")
		if err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
		}
	}
}
