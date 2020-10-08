package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (t *Telegram) commandUnrecognized(update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Unrecognized command"), nil
}

func (t *Telegram) showHelp(update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Help text"), nil
}
