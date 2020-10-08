package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (t *Telegram) commandNewSubnet(update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Choose subnet CIDR"), nil
}

func (t *Telegram) provideSubnetCIDR(update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(update.Message.Chat.ID, "Your subnet created"), nil
}
