package telegram

import (
	"github.com/cordalace/wireguard-for-homies/internal/db"
)

type DB interface {
	db.DB
	Begin(mode db.TxMode) (Tx, error)
}

type Tx interface {
	db.Tx
	CreateChatState(chatID int64, chatState ChatState) (ChatState, error)
	GetChatState(chatID int64) (ChatState, error)
	UpdateChatState(chatID int64, chatState ChatState) (ChatState, error)
	DeleteChatState(chatID int64) error
	CreateChatLongRunningCommand(chatID int64, command string) (string, error)
	GetChatLongRunningCommand(chatID int64) (string, error)
	UpdateChatLongRunningCommand(chatID int64, command string) (string, error)
	DeleteChatLongRunningCommand(chatID int64) error
}
