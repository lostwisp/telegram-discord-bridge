package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
	tgdis "github.com/lostwisp/telegram-discord-bridge/gRPC/tg-dis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	tgdis.UnimplementedTgdisMessageServer
	bot *discordgo.Session
	db  *pgx.Conn
}

func (s *server) NewMessage(ctx context.Context, req *tgdis.MessageRequest) (*tgdis.MessageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	chennelId, err := getchannelID(ctx, s.db, req.UserId)
	if err != nil {
		return &tgdis.MessageResponse{Score: 1}, status.Error(codes.NotFound, "User does not exist")
	}
	err = sendMessage(s.bot, chennelId, req.UserMessage)
	if err != nil {
		return &tgdis.MessageResponse{Score: 1}, status.Error(codes.Unknown, "The message is not sent")
	}
	return &tgdis.MessageResponse{Score: 0}, nil
}

func GetEnv() string {
	return os.Getenv("DISCORD_BOT_TOKEN")
}

func LoadConfig(TOKEN *string, host *string, port *string, db_name *string, user *string, password *string) {
	*TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
	*host = os.Getenv("HOST")
	*port = os.Getenv("PORT")
	*db_name = os.Getenv("DBNAME")
	*user = os.Getenv("USER")
	*password = os.Getenv("PASSWORD")
}

func main() {
	//Подключение к базе данных
	var TOKEN, host, port, db_name, user, password string
	LoadConfig(&TOKEN, &host, &port, &db_name, &user, &password)

	urldb := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, db_name)

	ctx, cancel := context.WithTimeout(context.Background(), (3 * time.Second))
	db, err := pgx.Connect(ctx, urldb)
	if err != nil {
		log.Println(err)
	}

	if db.Ping(ctx) != nil {
		log.Println("Ping error")
	}
	defer cancel()

	//Создание сервера для прослушивания gRPC
	conn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Println(err)
	}

	grpcServer := grpc.NewServer()
	server := &server{}
	tgdis.RegisterTgdisMessageServer(grpcServer, server)

	//Создание бота с WebSocket
	token := GetEnv()
	log.Println("Bot " + token)
	server.bot, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Println(err)
	}
	server.bot.AddHandler(pingPong)

	server.bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Ping, test the connect to the bot",
		})
	})

	err = server.bot.Open()
	if err != nil {
		log.Println(err)
	}
	defer server.bot.Close()

	//запуск gRPC
	go func() {
		if err := grpcServer.Serve(conn); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}

func getchannelID(ctx context.Context, db *pgx.Conn, user string) (string, error) {
	var channelId string
	err := db.QueryRow(ctx, "SELECT chennelId FROM user_channels WHERE telegramId=$1", user).Scan(&channelId)
	if err != nil {
		return "", err
	}
	return channelId, nil
}

func sendMessage(s *discordgo.Session, channelID string, content string) error {
	_, err := s.ChannelMessageSend(channelID, content)
	if err != nil {
		return err
	}
	return nil
}

func pingPong(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if m.Type == discordgo.InteractionApplicationCommand {
		log.Println(m.Type, "==", discordgo.InteractionApplicationCommand.String())
		s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong",
			},
		})
	}
}
