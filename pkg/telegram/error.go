package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errInvalidUrl        = errors.New("URL is invalid")
	errUserNotAuthorized = errors.New("User is not authorized")
	errUnabledToSave     = errors.New("Unabled to save")
)

func (b *Bot) handleError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, b.messages.Unknown)
	switch err {
	case errInvalidUrl:
		msg.Text = b.messages.InvalidUrl
	case errUserNotAuthorized:
		msg.Text = b.messages.Unauthorized
	case errUnabledToSave:
		msg.Text = b.messages.UnableToSave
	}

	b.bot.Send(msg)
}
