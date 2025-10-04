package bd

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func (*Storage) LoadConfig() (*config.ConfigDB, error) {
	return &config.ConfigDB{
		db_host: "DB_HOST",
	}
}

func (*Storage) CreateURL() {
	//postgres://[user]:[password]@[host]:[port]/[database]

	url, err := fmt.Sprintf("postgres://%s:%d@%s:%d/%s", config.ConfigDB)
}

//func (*Storage)

func (*Storage) ConnectDB() {
	db, err := pgxpool.New(context.Background())
}
