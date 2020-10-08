package badgerdb

import "github.com/cordalace/wireguard-for-homies/internal/telegram"

func (t *BadgerTx) GetUserChatState(userID int) (telegram.ChatState, error) {
	return telegram.ChatStateInitial, nil
}

func (t *BadgerTx) SetUserChatState(userID int, chatState telegram.ChatState) error {
	return nil
}
