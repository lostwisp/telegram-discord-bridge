package config

type ConfigDB struct {
	db_host    string
	db_port    string
	db_user    string
	db_pasword string
	db_name    string
}

var (
	TelegramToken string
)

var TelegramIdAdmin = map[int64]bool{
	1996660543: true,
	5526345699: true,
	5497536893: true,
}

var (
	DiscordToken     string
	DiscordchannelID string
)
