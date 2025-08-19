package discord

import (
	"github.com/HamsterNiki/TelegramDiscordBridge/MyDiscordPackege/config"
	"github.com/bwmarrin/discordgo"
)

var discord *discordgo.Session
var channelID string

func InstalDiscordBot() error {
	channelID = ""
	var err error
	discord, err = discordgo.New("Bot " + "")
	if err != nil {
		return err
	}
	err = discord.Open()
	if err != nil {
		print(err)
	}
	return nil
}

//func ReceivingMessagesAddHandler() {

//}

func SendMessengeToDiscord(massageTG string) {
	_, err := discord.ChannelMessageSend(channelID, massageTG)
	if err != nil {
		print(err)
	}
}
