package boltdb

import (
	"errors"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/mbakumenkov/go-pocket-bot/pkg/repository"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(chatId int64, token string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatId), []byte(token))
	})
}

func (r *TokenRepository) Get(chatId int64, bucket repository.Bucket) (string, error) {
	var token string
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		token = string(b.Get(intToBytes(chatId)))
		if token == "" {
			return errors.New("token not found")
		}
		return nil
	})
	return token, err
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
