package badgerdb

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
)

func fmtDBKey(prefix, id string) []byte {
	return []byte(fmt.Sprintf("%v:%v", prefix, id))
}

// getOrCreate returns value if key exists or returns defaultValue and creates key.
func (t *BadgerTx) getOrCreate(key string, defaultValue []byte) ([]byte, error) {
	var result []byte
	item, err := t.txn.Get([]byte(key))

	switch err {
	case badger.ErrKeyNotFound:
		errSet := t.txn.Set([]byte(key), defaultValue)
		if errSet != nil {
			return nil, errSet
		}
		result = append(result[:0], defaultValue...)

		return result, nil
	case nil:
		var errValueCopy error
		result, errValueCopy = item.ValueCopy(result)
		if errValueCopy != nil {
			return nil, errValueCopy
		}

		return result, nil
	default:
		return nil, err
	}
}

// getOrDefault returns value if key exists or returns defaultValue and doesn't writes anything to badger.
func (t *BadgerTx) getOrDefault(key string, defaultValue []byte) ([]byte, error) {
	var result []byte
	item, err := t.txn.Get([]byte(key))
	switch err {
	case badger.ErrKeyNotFound:
		result = append(result[:0], defaultValue...)

		return result, nil
	case nil:
		var errValueCopy error
		result, errValueCopy = item.ValueCopy(result)
		if errValueCopy != nil {
			return nil, errValueCopy
		}

		return result, nil
	default:
		return nil, err
	}
}

func (t *BadgerTx) exists(key []byte) (bool, error) {
	_, err := t.txn.Get(key)
	switch err {
	case badger.ErrKeyNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}
