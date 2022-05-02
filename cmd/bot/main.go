package main

import (
	"github.com/mbakumenkov/go-pocket-bot/pkg/config"
	"log"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mbakumenkov/go-pocket-bot/pkg/repository"
	"github.com/mbakumenkov/go-pocket-bot/pkg/repository/boltdb"
	"github.com/mbakumenkov/go-pocket-bot/pkg/server"
	"github.com/mbakumenkov/go-pocket-bot/pkg/telegram"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	bot.Debug = true
	if err != nil {
		log.Fatal(err)
	}

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)
	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages)
	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotURL)

	go func() {
		if err = telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err = authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens)); err != nil {
			return err
		}

		if _, err := tx.CreateBucketIfNotExists([]byte(repository.RequestTokens)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
