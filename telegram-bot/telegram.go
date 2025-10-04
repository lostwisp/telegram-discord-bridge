package telegram

import (
	"encoding/json"
	"fmt"

	"github.com/HamsterNiki/TelegramDiscordBridge/discord"
	"io"
	"net/http"
	"net/url"
	"time"

	config "github.com/HamsterNiki/TelegramDiscordBridge/pwd"
)

func main() {

}

func SendMessenge(idchat int, messege string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", config.TelegramToken, idchat, url.QueryEscape("Сообщение принято"))
	_, err := http.Get(url)
	discord.SendMessengeToDiscord(messege)
	if err != nil {
		println(err)
	}
	return nil
}

func Start(TOKEN string) error {
	offset := 791038172
	err := discord.InstalDiscordBot()
	if err != nil {
		println(err)
	}
	defer discord.Session.Close()

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
				Message  *struct {
					Text string `json:"text"`
					Chat struct {
						Id int `json:"id"`
					} `json:"chat"`
					From struct {
						ID int64 `json:"id"`
					} `json:"from"`
				} `json:"message"`
				CallbackQuery struct {
					ID   string `json:"id"`
					From struct {
						Id        int64  `json:"id"`
						IsBot     bool   `json:"is_bot"`
						FirstName string `json:"first_name"`
						Username  string `json:"username"`
					} `json:"from"`
					Data string `json:"data"`
				} `json:"callback_query"`
			} `json:"result"`
		}
		err = json.Unmarshal(dat, &result)
		if err != nil {
			println(err)
		}

		if result.Ok && (len(result.Result) > 0) {
			for _, update := range result.Result {

				println("UpdateID:", update.UpdateID, " UsersId:", update.Message.From.ID, " Message: ", update.Message.Text)
				if config.TelegramIdAdmin[update.Message.From.ID] {
					SendMessenge(update.Message.Chat.Id, update.Message.Text)
					if err != nil {
						println("Сообщение не отправленно")
					}
				} else {
					print("Доступ запрещён")
				}
				offset = update.UpdateID + 1
			}
		}

		time.Sleep(1 * time.Second)
	}
}
