package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (t *Telegram) commandUnrecognized(input *transitionInput) (*transitionOutput, error) {
	return &transitionOutput{
		responseMessage: tgbotapi.NewMessage(input.update.Message.Chat.ID, "Unrecognized command."),
		newState:        input.currentState,
	}, nil
}

func (t *Telegram) showHelp(input *transitionInput) (*transitionOutput, error) {
	return &transitionOutput{
		responseMessage: tgbotapi.NewMessage(input.update.Message.Chat.ID, "Help text"),
		newState:        input.currentState,
	}, nil
}
