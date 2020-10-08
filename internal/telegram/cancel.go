package telegram

import (
	"github.com/cordalace/wireguard-for-homies/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	responseNoActiveCommand = "No active command to cancel. I wasn't doing anything anyway. Zzzzz..."
	responseCommandCanceled = "Command has been cancelled. Anything else I can do for you?"
)

func (t *Telegram) commandCancel(update *tgbotapi.Update) (tgbotapi.Chattable, error) {
	tx, err := t.db.Begin(db.TxModeReadWrite)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	state, err := tx.GetUserChatState(update.Message.From.ID)
	if err != nil {
		return nil, err
	}

	if state == ChatStateInitial {
		return tgbotapi.NewMessage(update.Message.Chat.ID, responseNoActiveCommand), nil
	}

	err = tx.SetUserChatState(update.Message.From.ID, ChatStateInitial)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, responseCommandCanceled), nil
}
