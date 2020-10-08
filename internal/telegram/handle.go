package telegram

import (
	"errors"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

var errUnknownState = errors.New("unknown state")

func (t *Telegram) handleUpdate(update *tgbotapi.Update) error {
	if update.Message == nil { // ignore any non-Message Updates
		return nil
	}

	t.logger.Info(
		"telegram message received",
		zap.String("userName", update.Message.From.UserName),
		zap.String("text", update.Message.Text),
	)

	var handler func(*tgbotapi.Update) (tgbotapi.Chattable, error)

	tx, err := t.db.Begin(db.TxModeReadOnly)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	state, err := tx.GetUserChatState(update.Message.From.ID)
	if err != nil {
		return err
	}

	if update.Message.IsCommand() && update.Message.Command() == "/cancel" {
		handler = t.commandCancel
	} else {
		switch state {
		case ChatStateInitial:
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "/newsubnet":
					handler = t.commandNewSubnet
				default:
					handler = t.commandUnrecognized
				}
			} else {
				handler = t.showHelp
			}
		case ChatStateSubnetExpectCIDR:
			handler = t.provideSubnetCIDR
		default:
			return errUnknownState
		}
	}

	responseMessage, err := handler(update)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, error happened")
		if _, sorrySendErr := t.bot.Send(msg); sorrySendErr != nil {
			return sorrySendErr
		}
	}

	if _, err = t.bot.Send(responseMessage); err != nil {
		return err
	}

	return nil
}
