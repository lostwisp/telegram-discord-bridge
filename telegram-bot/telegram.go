package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	tgdis "github.com/lostwisp/telegram-discord-bridge/gRPC/tg-dis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var TelegramIdAdmin = map[int64]bool{
	1996660543: true,
	5526345699: true,
	5497536893: true,
}

type ResultJSON struct {
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

type DiscordgRPC struct {
	Client tgdis.TgdisMessageClient
	conn   *grpc.ClientConn
}

func NewDiscordgRPC() *DiscordgRPC {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}

	Client := tgdis.NewTgdisMessageClient(conn)
	return &DiscordgRPC{Client, conn}
}

func (D *DiscordgRPC) SendToDiscord(mes string) (*tgdis.MessageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req := &tgdis.MessageRequest{Message: mes}
	resp, err := D.Client.NewMessage(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return resp, nil
}

type TelegramBot struct {
	URL    string
	Offset int
}

type bdconfig struct {
	host     string
	port     string
	db_name  string
	user     string
	password string
}

func (T *TelegramBot) Update() (*ResultJSON, error) {

	resp, err := http.Get(T.URL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp.Body.Close()

	var result ResultJSON
	err = json.Unmarshal(dat, &result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if result.Ok && (len(result.Result) > 0) {
		for _, update := range result.Result {
			T.Offset = update.UpdateID + 1
		}
	}

	return &result, nil
}

func (T TelegramBot) SendToUser(idchat int, mes string) (*http.Response, error) {
	url := fmt.Sprintf("%s/sendMessage", T.URL)
	var str = struct {
		ChatID  int    `json:"chat_id"`
		Message string `json:"text"`
	}{
		ChatID:  idchat,
		Message: mes,
	}
	bytestr, err := json.Marshal(str)
	if err != nil {
		log.Println(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytestr))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (T TelegramBot) Requestbd() error {
	return nil
}

func (T TelegramBot) Save() error {
	return nil
}
func LoadConfig(TOKEN *string, host *string, port *string, db_name *string, user *string, password *string) {
	*TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
	*host = os.Getenv("HOST")
	*port = os.Getenv("PORT")
	*db_name = os.Getenv("DBNAME")
	*user = os.Getenv("USER")
	*password = os.Getenv("PASSWORD")
}

func main() {

	var TOKEN, host, port, db_name, user, password string
	LoadConfig(&TOKEN, &host, &port, &db_name, &user, &password)

	urldb := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, db_name)

	ctx, cancel := context.WithTimeout(context.Background(), (time.Second * 3))

	conn, err := pgx.Connect(ctx, urldb)
	if err != nil {
		log.Println(err)
	}

	if conn != nil {
		err = conn.Ping(ctx)
		if err != nil {
			log.Println("Failed to connect", err)
		}
	}

	defer cancel()
	url := fmt.Sprintf("https://api.telegram.org/bot%s", TOKEN)
	bot := TelegramBot{url, 791038172}

	clintDiscord := NewDiscordgRPC()
	defer clintDiscord.conn.Close()

	for {

		time.Sleep(1 * time.Second)
		result, err := bot.Update()
		if err != nil {
			log.Println(err)
		}

		for _, update := range result.Result {
			resp, err := bot.SendToUser(update.Message.Chat.Id, update.Message.Text)
			if err != nil {
				log.Println(err)
			}
			if resp.StatusCode != 200 {
				log.Println("The message was not delivered", resp.Status)
			}
			mesResp, err := clintDiscord.SendToDiscord(update.Message.Text)
			if err != nil {
				log.Println(err)
			}
			if mesResp.Score != 0 {

				log.Println("Services Discord didn't send the message.")
			}

		}
	}
}

func ConnectDB() {

}
