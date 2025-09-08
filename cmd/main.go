package main

import (
	"fmt"
	"github.com/HamsterNiki/TelegramDiscordBridge/pwd"
	"github.com/HamsterNiki/TelegramDiscordBridge/telegram"
)

func main() {
	var err error
	fmt.Print("Telegram token: ")
	_, err = fmt.Scanln(&config.TelegramToken)
	if err != nil {
		println(err)
	}
	fmt.Print("Discord token: ")
	_, err = fmt.Scanln(&config.DiscordToken)
	if err != nil {
		println(err)
	}
	telegram.Start(config.TelegramToken)
}
