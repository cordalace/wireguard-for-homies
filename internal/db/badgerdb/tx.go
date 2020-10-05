package badgerdb

import badger "github.com/dgraph-io/badger/v2"

type BadgerTx struct {
	txn *badger.Txn
}

func (t *BadgerTx) Commit() error {
	return t.txn.Commit()
}

func (t *BadgerTx) Rollback() {
	t.txn.Discard()
}
