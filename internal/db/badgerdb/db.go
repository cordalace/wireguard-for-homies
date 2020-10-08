package badgerdb

import (
	"errors"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/cordalace/wireguard-for-homies/internal/manager"
	"github.com/cordalace/wireguard-for-homies/internal/telegram"
	badger "github.com/dgraph-io/badger/v2"
)

var errNotInitialized = errors.New("badger db is not initialized, call Init() first")

type BadgerDB struct {
	db   *badger.DB
	opts badger.Options
}

type badgerManagerDB struct {
	BadgerDB
}

type badgerTelegramDB struct {
	BadgerDB
}

func NewBadgerDB(opts badger.Options) *BadgerDB {
	return &BadgerDB{opts: opts}
}

func (d *BadgerDB) AsManagerDB() manager.DB {
	return &badgerManagerDB{BadgerDB{db: d.db, opts: d.opts}}
}

func (d *BadgerDB) AsTelegramDB() telegram.DB {
	return &badgerTelegramDB{BadgerDB{db: d.db, opts: d.opts}}
}

func (d *BadgerDB) Init() error {
	var err error
	d.db, err = badger.Open(d.opts)

	return err
}

func (d *BadgerDB) Close() error {
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

func (d *BadgerDB) Begin(mode db.TxMode) (*BadgerTx, error) {
	if d.db == nil {
		return nil, errNotInitialized
	}
	update := mode == db.TxModeReadWrite
	txn := d.db.NewTransaction(update)

	return &BadgerTx{txn: txn}, nil
}

func (d *badgerManagerDB) Begin(mode db.TxMode) (manager.Tx, error) {
	return d.BadgerDB.Begin(mode)
}

func (d *badgerTelegramDB) Begin(mode db.TxMode) (telegram.Tx, error) {
	return d.BadgerDB.Begin(mode)
}
