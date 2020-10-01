package badgerdb

import badger "github.com/dgraph-io/badger/v2"

type badgerTx struct {
	txn *badger.Txn
}

func (t *badgerTx) Commit() error {
	return t.txn.Commit()
}

func (t *badgerTx) Rollback() {
	t.txn.Discard()
}
