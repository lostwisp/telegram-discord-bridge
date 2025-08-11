package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var TOKEN string

func SendMessenge(idchat int, messege string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", TOKEN, idchat, url.QueryEscape(messege))
	_, err := http.Get(url)
	if err != nil {
		println(err)
	}
	return nil
}

func main() {
	offset := 791038172
	_, err := fmt.Scanln(&TOKEN)
	if err != nil {
		println(err)
	}
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
				SendMessenge(update.Message.Chat.Id, "Сообщение принято")
				if err != nil {
					println("Сообщение не отправленно")
				}
				offset = update.UpdateID + 1
			}
		}

		time.Sleep(1 * time.Second)
	}
}
