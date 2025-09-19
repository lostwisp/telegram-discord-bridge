package discord

import (
	config "github.com/HamsterNiki/TelegramDiscordBridge/pwd"
	"github.com/bwmarrin/discordgo"
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
