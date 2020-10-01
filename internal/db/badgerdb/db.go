package badgerdb

import (
	"errors"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	badger "github.com/dgraph-io/badger/v2"
)

var errNotInitialized = errors.New("badger db is not initialized, call Init() first")

type badgerDB struct {
	db   *badger.DB
	opts badger.Options
}

func NewBadgerDB(opts badger.Options) db.DB {
	return &badgerDB{opts: opts}
}

func (d *badgerDB) Init() error {
	var err error
	d.db, err = badger.Open(d.opts)

	return err
}

func (d *badgerDB) Close() error {
	if d.db == nil {
		return nil
	}
	err := d.db.Close()
	if err != nil {
		return err
	}
	d.db = nil
	return nil
}

func (d *badgerDB) Begin(mode db.TxMode) (db.Tx, error) {
	if d.db == nil {
		return nil, errNotInitialized
	}
	update := mode == db.TxModeReadWrite
	txn := d.db.NewTransaction(update)

	return &badgerTx{txn: txn}, nil
}
