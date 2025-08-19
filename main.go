package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"tg-discord-bot/MyDiscord"
	"tg-discord-bot/MyTelegram"
	"time"
)

var TOKEN string

func SendMessenge(idchat int, messege string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", TOKEN, idchat, url.QueryEscape("Сообщение принято"))
	_, err := http.Get(url)
	MyDiscord.SendMessengeToDiscord(messege)
	if err != nil {
		println(err)
	}
	return nil
}

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

	go MyTelegram.

}
