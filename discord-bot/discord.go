package discord

import (
	"github.com/bwmarrin/discordgo"
	tgdis "github.com/lostwisp/telegram-discord-bridge/gRPC/tg-dis"
	"log"
	"net"
)

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
}
