package discord

import (
	"context"
	"log"
	"net"

	"github.com/bwmarrin/discordgo"
	tgdis "github.com/lostwisp/telegram-discord-bridge/gRPC/tg-dis"
	config "github.com/lostwisp/telegram-discord-bridge/pwd"
	"google.golang.org/grpc"
)

type server struct {
	tgdis.UnimplementedTgdisMessageServiceServer
}

func (s *server) NewMessage(ctx context.Context, req *tgdis.MessageRequest) (*tgdis.MessageResponse, error) {
	return &tgdis.MessageResponse{Score: 0}, nil
}

var (
	DiscordToken     string
	DiscordchannelID string
)

var Session *discordgo.Session

func InstalDiscordBot() error {
	var err error
	Session, err = discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return err
	}
	err = Session.Open()
	if err != nil {
		print(err)
	}
	return nil

}

func SendMessengeToDiscord(massageTG string) {
	_, err := Session.ChannelMessageSend(config.DiscordchannelID, massageTG)
	if err != nil {
		print(err)
	}
}

func main() {
	conn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Println(err)
	}

	grpcServer := grpc.NewServer()

	tgdis.RegisterTgdisMessageServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
