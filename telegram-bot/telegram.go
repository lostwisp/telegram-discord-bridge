package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	tgdis "github.com/lostwisp/telegram-discord-bridge/gRPC/tg-dis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Result struct {
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
}

type ResultJSON struct {
	Ok     bool     `json:"ok"`
	Result []Result `json:"result"`
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
		return nil, err
	}
	return resp, nil
}

type command int

const (
	commandSetChannell command = iota
	commandPing
)

type TelegramBot struct {
	URL        string
	Offset     int
	commandNow map[string]command
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
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytestr))
	if err != nil {
		return nil, err
	}
	return resp, nil
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
	bot := TelegramBot{url, 791038172, make(map[string]command)}

	clintDiscord := NewDiscordgRPC()
	defer clintDiscord.conn.Close()

	for {

		time.Sleep(1 * time.Second)
		result, err := bot.Update()
		if err != nil {
			log.Println(err)
		}

		for _, update := range result.Result {
			go process(&update)
			//Добавть привязку чата к акаунту телеграмм
			ctx, cancel := context.WithTimeout(context.Background(), (time.Second * 3))
			defer cancel()
			var chatId = getDiscordChat(ctx, conn, string(update.Message.From.ID))
			if chatId == "" {
				bot.SendToUser(update.Message.Chat.Id, "Привяжите свой акаунт телеграм к дискорд чату")
			}

			//В дискорд
			mesResp, err := clintDiscord.SendToDiscord(update.Message.Text)
			if err != nil {
				log.Println(err)
			}
			if mesResp.Score == 0 {
				log.Printf("Discord service delivered the message from user %d\n", update.Message.From.ID)
				resp, err := bot.SendToUser(update.Message.Chat.Id, "Сообщение отправленно")
				if err != nil {
					log.Println(err)
				}
				if resp.StatusCode != 200 {
					log.Println("The message was not delivered to user in telegram chat", resp.Status)
				}
			} else {
				log.Printf("Discord service don't deliver the message from user %d\n", update.Message.From.ID)
				resp, err := bot.SendToUser(update.Message.Chat.Id, "Ошибка, сообщение не отправленно")
				if err != nil {
					log.Println(err)
				}
				if resp.StatusCode != 200 {
					log.Println("The message was not delivered to user in telegram chat", resp.Status)
				}
			}

		}
	}
}

func process(update *Result) {

}

func InitTables(ctx context.Context, conn *pgx.Conn) (pgconn.CommandTag, error) {
	tag, err := conn.Exec(ctx, "CREATE TABLE IF NOT EXIST user_channels {telegramId TEXT PRIMARY KEY, chennelId TEXT}")
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	return tag, nil
}

func getDiscordChat(ctx context.Context, conn *pgx.Conn, userID string) string {
	var channelId string
	err := conn.QueryRow(ctx, "SELECT chennelId FROM user_channels WHERE telegramId=$1", userID).Scan(&channelId)
	if errors.Is(err, pgx.ErrNoRows) {
		return ""
	}
	return channelId
}

func setDiscordChat(ctx context.Context, conn *pgx.Conn, userID string, channelId string) error {
	_, err := conn.Exec(ctx, "INSERT INTO user_channels VALUES (%1, %2) ON CONFLICT(telegramId) DO UPDATE SET chennelId=EXCLUDED(user_channels)", userID, channelId)
	if err != nil {
		return err
	}
	return nil
}
