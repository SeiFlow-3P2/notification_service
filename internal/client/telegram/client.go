package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
)

type Client struct {
	bot *tgbotapi.BotAPI
}

func NewClient(token string) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Client{bot: bot}, nil
}

func (c *Client) SendMessage(userID int64, text string) error {
	msg := tgbotapi.NewMessage(userID, text)
	_, err := c.bot.Send(msg)
	return err
}