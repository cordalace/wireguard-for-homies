package badgerdb

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

func getOrCreate(txn *badger.Txn, key string, value []byte) ([]byte, error) {
	var result []byte
	item, err := txn.Get([]byte(key))

	switch err {
	case badger.ErrKeyNotFound:
		errSet := txn.Set([]byte(key), value)
		if errSet != nil {
			return nil, errSet
		}
		result = append(result[:0], value...)

		return result, nil
	case nil:
		var copyErr error
		result, copyErr = item.ValueCopy(result)

		if copyErr != nil {
			return nil, copyErr
		}

		return result, nil
	default:
		return nil, err
	}
}

func fmtDBKey(prefix, id string) []byte {
	return []byte(fmt.Sprintf("%v:%v", prefix, id))
}
