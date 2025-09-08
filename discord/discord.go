package discord

import (
	"github.com/bwmarrin/discordgo"
)

var Session *discordgo.Session
var channelID string

func InstalDiscordBot() error {
	channelID = ""
	var err error
	Session, err = discordgo.New("Bot " + "")
	if err != nil {
		return err
	}
	err = Session.Open()
	if err != nil {
		print(err)
	}
	return nil
}

//func ReceivingMessagesAddHandler() {

//}

func SendMessengeToDiscord(massageTG string) {
	_, err := Session.ChannelMessageSend(channelID, massageTG)
	if err != nil {
		print(err)
	}
}
