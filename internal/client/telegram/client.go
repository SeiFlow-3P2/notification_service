package telegram

import (
    "context"
    "errors"
    "strings"
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

func (c *Client) SendMessage(ctx context.Context, userID int64, text string) error {
    msg := tgbotapi.NewMessage(userID, text)
    _, err := c.bot.Send(msg)
    if err != nil {
        errStr := err.Error()
        if strings.Contains(errStr, "403") {
            return errors.New("bot is blocked by user or user not found")
        } else if strings.Contains(errStr, "429") {
            return errors.New("too many requests to Telegram API")
        }
        return err
    }
    return nil
}