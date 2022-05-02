package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mbakumenkov/go-pocket-bot/pkg/repository"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(b.messages.Start, authLink))
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) getAccessToken(chatId int64) (string, error) {
	return b.tokenRepository.Get(chatId, repository.AccessTokens)
}

func (b *Bot) generateAuthorizationLink(chatId int64) (string, error) {
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectUrl)
	if err != nil {
		return "", err
	}

	b.tokenRepository.Save(chatId, requestToken, repository.RequestTokens)

	return b.pocketClient.GetAuthorizationURL(requestToken, b.generateRedirectLink(chatId))
}

func (b *Bot) generateRedirectLink(chatId int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatId)
}
