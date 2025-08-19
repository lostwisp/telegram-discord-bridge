package MyTelegram

import (
	"encoding/json"
	"fmt"
	"github.com/HamsterNiki/TelegramDiscordBridge/MyDiscord"
	"github.com/HamsterNiki/TelegramDiscordBridge/main"
	"io"
	"net/http"
	"net/url"
	"time"
)

func SendMessenge(idchat int, messege string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", telegramToken, idchat, url.QueryEscape("Сообщение принято"))
	_, err := http.Get(url)
	MyDiscord.SendMessengeToDiscord(messege)
	if err != nil {
		println(err)
	}
	return nil
}

func Start(TOCEN string) error {
	offset := 791038172
	_, err := fmt.Scanln(&TOKEN)
	if err != nil {
		println(err)
	}

	err = InstalDiscordBot()
	if err != nil {
		println(err)
	}
	defer discord.Close()

	for {
		URL := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", TOKEN, offset)
		resp, err := http.Get(URL)
		if err != nil {
			println(err)
		}
		dat, err := io.ReadAll(resp.Body)
		if err != nil {
			println(err)
		}
		resp.Body.Close()
		var result struct {
			Ok     bool `json:"ok"`
			Result []struct {
				UpdateID int `json:"update_id"`
				Message  struct {
					Text string `json:"text"`
					Chat struct {
						Id int `json:"id"`
					} `json:"chat"`
				} `json:"message"`
			} `json:"result"`
		}
		err = json.Unmarshal(dat, &result)
		if err != nil {
			println(err)
		}

		if result.Ok && (len(result.Result) > 0) {
			for _, update := range result.Result {
				println(update.UpdateID, " ", update.Message.Text)
				SendMessenge(update.Message.Chat.Id, update.Message.Text)
				if err != nil {
					println("Сообщение не отправленно")
				}
				offset = update.UpdateID + 1
			}
		}

		time.Sleep(1 * time.Second)
	}
}
