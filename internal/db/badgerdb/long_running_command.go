package badgerdb

import (
	"encoding/json"
	"strconv"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	badger "github.com/dgraph-io/badger/v2"
)

const longRunningCommandPrefix = "chatLongRunningCommand"

type longRunningCommand string

func (c longRunningCommand) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

func longRunningCommandFromJSON(data []byte) (string, error) {
	var ret string

	err := json.Unmarshal(data, &ret)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func (t *BadgerTx) CreateChatLongRunningCommand(chatID int64, command string) (string, error) {
	key := fmtDBKey(longRunningCommandPrefix, strconv.FormatInt(chatID, 10))

	exists, err := t.exists(key)
	if err != nil {
		return "", err
	}
	if exists {
		return "", db.ErrAlreadyExists
	}

	commandJSON, err := longRunningCommand(command).ToJSON()
	if err != nil {
		return "", err
	}
	err = t.txn.Set(key, commandJSON)
	if err != nil {
		return "", err
	}

	return command, nil
}

func (t *BadgerTx) GetChatLongRunningCommand(chatID int64) (string, error) {
	key := fmtDBKey(longRunningCommandPrefix, strconv.FormatInt(chatID, 10))
	item, err := t.txn.Get(key)
	switch err {
	case badger.ErrKeyNotFound:
		return "", db.ErrNotFound
	case nil:
		var (
			value   []byte
			copyErr error
		)
		value, copyErr = item.ValueCopy(value)
		if copyErr != nil {
			return "", copyErr
		}
		return longRunningCommandFromJSON(value)
	default:
		return "", err
	}
}

func (t *BadgerTx) UpdateChatLongRunningCommand(chatID int64, command string) (string, error) {
	key := fmtDBKey(longRunningCommandPrefix, strconv.FormatInt(chatID, 10))

	commandJSON, err := longRunningCommand(command).ToJSON()
	if err != nil {
		return "", err
	}
	err = t.txn.Set(key, commandJSON)
	if err != nil {
		return "", err
	}

	return command, nil
}

func (t *BadgerTx) DeleteChatLongRunningCommand(chatID int64) error {
	key := fmtDBKey(longRunningCommandPrefix, strconv.FormatInt(chatID, 10))

	exists, err := t.exists(key)
	if err != nil {
		return err
	}
	if !exists {
		return db.ErrNotFound
	}

	return t.txn.Delete(key)
}
