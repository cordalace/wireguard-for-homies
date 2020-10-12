package telegram

import (
	"errors"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type messageType int

const (
	messageTypeCmdUnknown messageType = iota
	messageTypeText
	messageTypeCmdCancel
	messageTypeCmdHelp
	messageTypeCmdNewSubnet
)

type fsmKey struct {
	currentState ChatState
	event        messageType
}

type handlerFunc func(input *transitionInput) (*transitionOutput, error)

type fsm map[fsmKey]handlerFunc

var errUnknownTransition = errors.New("unknown transition")

type transitionInput struct {
	update       *tgbotapi.Update
	currentState ChatState
	txReadWrite  Tx
	txReadOnly   Tx
}

type transitionOutput struct {
	responseMessage tgbotapi.Chattable
	newState        ChatState
}

func (t *Telegram) getMessageType(msg *tgbotapi.Message) messageType {
	if msg.IsCommand() {
		msgType, ok := t.messageTypeByCommand[msg.Command()]
		if !ok {
			return messageTypeCmdUnknown
		}
		return msgType
	}
	return messageTypeText
}

func (t *Telegram) withReadWriteTx(handler handlerFunc) handlerFunc {
	return func(input *transitionInput) (*transitionOutput, error) {
		var err error
		input.txReadWrite, err = t.db.Begin(db.TxModeReadWrite)
		if err != nil {
			return nil, err
		}
		defer input.txReadWrite.Rollback()

		output, err := handler(input)
		if err != nil {
			return nil, err
		}

		err = input.txReadWrite.Commit()
		if err != nil {
			return nil, err
		}

		return output, err
	}
}

// func (t *Telegram) withReadOnlyTx(handler handlerFunc) handlerFunc {
// 	return func(input *transitionInput) (*transitionOutput, error) {
// 		var err error
// 		input.txReadOnly, err = t.db.Begin(db.TxModeReadOnly)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer input.txReadOnly.Rollback()

// 		return handler(input)
// 	}
// }

func (t *Telegram) getHandler(update *tgbotapi.Update, currentState ChatState) (handlerFunc, error) {
	key := fsmKey{currentState: currentState, event: t.getMessageType(update.Message)}
	handler, ok := t.fsm[key]
	if !ok {
		return nil, errUnknownTransition
	}
	t.logger.Info("choosed handler", zap.Int("currentState", int(currentState)))
	return handler, nil
}

func (t *Telegram) handleUpdate(update *tgbotapi.Update) error {
	if update.Message == nil { // ignore any non-Message Updates
		return nil
	}

	t.logger.Info(
		"telegram message received",
		zap.String("userName", update.Message.From.UserName),
		zap.String("text", update.Message.Text),
	)

	txReadState, err := t.db.Begin(db.TxModeReadOnly)
	if err != nil {
		return err
	}
	defer txReadState.Rollback()

	currentState, err := txReadState.GetChatState(update.Message.Chat.ID)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return err
		}
		currentState = ChatStateInitial
	}

	handler, err := t.getHandler(update, currentState)
	if err != nil {
		return err
	}

	input := &transitionInput{
		update:       update,
		currentState: currentState,
		txReadWrite:  nil,
		txReadOnly:   nil,
	}
	output, err := handler(input)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, error happened")
		if _, sorrySendErr := t.bot.Send(msg); sorrySendErr != nil {
			return sorrySendErr
		}
		return err
	}

	if input.currentState != output.newState {
		txWriteState, err := t.db.Begin(db.TxModeReadWrite)
		if err != nil {
			return err
		}
		defer txWriteState.Rollback()

		if _, err := txWriteState.UpdateChatState(update.Message.Chat.ID, output.newState); err != nil {
			return err
		}

		if err := txWriteState.Commit(); err != nil {
			return err
		}

		t.logger.Info("state written", zap.Int64("chatID", update.Message.Chat.ID), zap.Int("state", int(output.newState)))
	}

	if _, err = t.bot.Send(output.responseMessage); err != nil {
		return err
	}

	return nil
}
