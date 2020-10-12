package telegram

import (
	"errors"
	"fmt"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	responseNoActiveCommand = "No active command to cancel. I wasn't doing anything anyway. Zzzzz..."
)

func (t *Telegram) commandCancel(input *transitionInput) (*transitionOutput, error) {
	if input.currentState == ChatStateInitial {
		return &transitionOutput{
			responseMessage: tgbotapi.NewMessage(input.update.Message.Chat.ID, responseNoActiveCommand),
			newState:        ChatStateInitial,
		}, nil
	}

	command, err := input.txReadWrite.GetChatLongRunningCommand(input.update.Message.Chat.ID)
	if err != nil {
		return nil, err
	}

	if err := input.txReadWrite.DeleteChatLongRunningCommand(input.update.Message.Chat.ID); err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}
	}

	responseText := fmt.Sprintf("The command %v has been cancelled.", command)
	return &transitionOutput{
		responseMessage: tgbotapi.NewMessage(input.update.Message.Chat.ID, responseText),
		newState:        ChatStateInitial,
	}, nil
}
