package telegram

import (
	"errors"
	"fmt"
	"net"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *Telegram) commandNewSubnet(input *transitionInput) (*transitionOutput, error) {
	if _, err := input.txReadWrite.UpdateChatLongRunningCommand(
		input.update.Message.Chat.ID,
		input.update.Message.Command(),
	); err != nil {
		return nil, err
	}

	return &transitionOutput{
		responseMessage: tgbotapi.NewMessage(input.update.Message.Chat.ID, "Choose subnet CIDR"),
		newState:        ChatStateSubnetExpectCIDR,
	}, nil
}

func (t *Telegram) provideSubnetCIDR(input *transitionInput) (*transitionOutput, error) {
	if _, _, err := net.ParseCIDR(input.update.Message.Text); err != nil {
		return &transitionOutput{
			responseMessage: tgbotapi.NewMessage(input.update.Message.Chat.ID, fmt.Sprintf("Invalid subnet CIDR: %v", err)),
			newState:        input.currentState,
		}, nil
	}

	// // should we normalize cidr?
	// subnet, err := t.wgManager.CreateSubnet(&models.Subnet{CIDR: ipNet})
	// if err != nil {
	// 	return nil, err
	// }

	if err := input.txReadWrite.DeleteChatLongRunningCommand(input.update.Message.Chat.ID); err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return nil, err
		}
	}

	return &transitionOutput{
		responseMessage: tgbotapi.NewMessage(input.update.Message.Chat.ID, "Your subnet created."),
		newState:        ChatStateInitial,
	}, nil
}
