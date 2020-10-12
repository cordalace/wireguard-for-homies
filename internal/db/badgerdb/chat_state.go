package badgerdb

import (
	"strconv"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/cordalace/wireguard-for-homies/internal/telegram"
	badger "github.com/dgraph-io/badger/v2"
)

const chatStatePrefix = "chatState"

func (t *BadgerTx) CreateChatState(chatID int64, chatState telegram.ChatState) (telegram.ChatState, error) {
	key := fmtDBKey(chatStatePrefix, strconv.FormatInt(chatID, 10))

	exists, err := t.exists(key)
	if err != nil {
		return telegram.ChatStateInitial, err
	}
	if exists {
		return telegram.ChatStateInitial, db.ErrAlreadyExists
	}

	chatStateJSON, err := chatState.ToJSON()
	if err != nil {
		return telegram.ChatStateInitial, err
	}
	err = t.txn.Set(key, chatStateJSON)
	if err != nil {
		return telegram.ChatStateInitial, err
	}

	return chatState, nil
}

func (t *BadgerTx) GetChatState(chatID int64) (telegram.ChatState, error) {
	key := fmtDBKey(chatStatePrefix, strconv.FormatInt(chatID, 10))
	item, err := t.txn.Get(key)
	switch err {
	case badger.ErrKeyNotFound:
		return telegram.ChatStateInitial, db.ErrNotFound
	case nil:
		var (
			value   []byte
			copyErr error
		)
		value, copyErr = item.ValueCopy(value)
		if copyErr != nil {
			return telegram.ChatStateInitial, copyErr
		}
		return telegram.ChatStateFromJSON(value)
	default:
		return telegram.ChatStateInitial, err
	}
}

func (t *BadgerTx) UpdateChatState(chatID int64, chatState telegram.ChatState) (telegram.ChatState, error) {
	key := fmtDBKey(chatStatePrefix, strconv.FormatInt(chatID, 10))

	chatStateJSON, err := chatState.ToJSON()
	if err != nil {
		return telegram.ChatStateInitial, err
	}
	err = t.txn.Set(key, chatStateJSON)
	if err != nil {
		return telegram.ChatStateInitial, err
	}

	return chatState, nil
}

func (t *BadgerTx) DeleteChatState(chatID int64) error {
	key := fmtDBKey(chatStatePrefix, strconv.FormatInt(chatID, 10))

	exists, err := t.exists(key)
	if err != nil {
		return err
	}
	if !exists {
		return db.ErrNotFound
	}

	return t.txn.Delete(key)
}
