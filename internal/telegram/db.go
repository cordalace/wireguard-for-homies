package telegram

import (
	"github.com/cordalace/wireguard-for-homies/internal/db"
)

type ChatState int

const (
	ChatStateInitial = iota
	ChatStateSubnetExpectCIDR
)

type DB interface {
	db.DB
	Begin(mode db.TxMode) (Tx, error)
}

type Tx interface {
	db.Tx
	GetUserChatState(userID int) (ChatState, error)
	SetUserChatState(userID int, chatState ChatState) error
}
