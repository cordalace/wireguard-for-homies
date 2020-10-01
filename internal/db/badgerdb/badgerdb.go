package badgerdb

import (
	"github.com/cordalace/wireguard-for-homies/internal/db"
	badger "github.com/dgraph-io/badger/v2"
)

type badgerDB struct {
	db   *badger.DB
	opts badger.Options
}

type badgerTX struct {
	txn *badger.Txn
}

func NewBadgerDB(opts badger.Options) db.DB {
	return &badgerDB{opts: opts}
}

func (d *badgerDB) Init() error {
	var err error
	d.db, err = badger.Open(d.opts)

	return err
}

func (d *badgerDB) Begin() (db.Tx, error) {
	txn := d.db.NewTransaction(true)

	return &badgerTX{txn: txn}, nil
}

func (d *badgerDB) Close() error {
	return d.db.Close()
}

func (t *badgerTX) Commit() error {
	return t.txn.Commit()
}

func (t *badgerTX) Rollback() {
	t.txn.Discard()
}
