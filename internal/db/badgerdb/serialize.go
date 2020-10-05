package badgerdb

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v2"
)

func (t *BadgerTx) DumpData() ([]byte, error) {
	dbContents := make(map[string]json.RawMessage)
	opts := badger.DefaultIteratorOptions
	// opts.PrefetchSize = 10
	it := t.txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		key := string(item.Key())
		err := item.Value(func(val []byte) error {
			dbContents[key] = val
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return json.MarshalIndent(dbContents, "", "    ")
}

func (t *BadgerTx) LoadData(data []byte) error {
	dbContents := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &dbContents)
	if err != nil {
		return err
	}

	for key, val := range dbContents {
		err = t.txn.Set([]byte(key), val)
		if err != nil {
			return err
		}
	}
	return nil
}
