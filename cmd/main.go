package main

import (
	"fmt"
	"github.com/HamsterWiki/TelegramDiscordBridge/discord"
	"github.com/HamsterWiki/TelegramDiscordBridge/telegram"
)

func main() {
	var err error
	fmt.Print("Telegram token: ")
	_, err = fmt.Scanln(&telegramToken)
	if err != nil {
		println(err)
	}

	fmt.Print("Discord token: ")
	_, err = fmt.Scanln(&discordToken)
	if err != nil {
		println(err)
	}

}
