package badgerdb

import (
	"bytes"
	"errors"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	badger "github.com/dgraph-io/badger/v2"
)

type initDB int

const (
	initDBEmpty = iota
	initDBWithInput
)

type txMode int

const (
	txModeReadOnly = iota
	txModeReadWrite
)

func withTestTx(t *testing.T, init initDB, mode txMode, testFunc func(txn *badger.Txn)) {
	var ddb *badger.DB
	switch init {
	case initDBEmpty:
		ddb = openInMemoryDB(t)
	case initDBWithInput:
		ddb = openInMemoryDBWithData(t)
	default:
		t.Fatalf("unknown init db: %v", init)
	}

	var txn *badger.Txn
	switch mode {
	case txModeReadOnly:
		txn = ddb.NewTransaction(false)
	case txModeReadWrite:
		txn = ddb.NewTransaction(true)
	default:
		t.Fatalf("unknown tx mode: %v", mode)
	}
	defer txn.Discard()

	testFunc(txn)

	if mode == txModeReadWrite {
		cupaloy.New(cupaloy.SnapshotFileExtension(".json")).SnapshotT(t, dumpData(t, &BadgerTx{txn: txn}))
	}
}

func setKeyValue(t *testing.T, txn *badger.Txn, key string, value []byte) {
	err := txn.Set([]byte(key), value)
	if err != nil {
		t.Fatalf("badger.Txn.Set() error = %v, want nil", err)
	}
}

func ensureKeyNotFound(t *testing.T, txn *badger.Txn, key string) {
	_, err := txn.Get([]byte(key))
	if !errors.Is(err, badger.ErrKeyNotFound) {
		t.Fatalf("badger.Txn.Get() error = %v, want Key not found", err)
	}
}

func assertKeyValue(t *testing.T, txn *badger.Txn, key string, wantValue []byte) {
	item, err := txn.Get([]byte(key))
	if err != nil {
		t.Fatalf("badger.Txn.Get() error = %v, want nil", err)
	}

	var gotValue []byte
	gotValue, err = item.ValueCopy(gotValue)
	if err != nil {
		t.Fatalf("badger.Item.ValueCopy() error = %v, want nil", err)
	}

	if !bytes.Equal(gotValue, wantValue) {
		t.Fatalf("badger.Txn.Get() = %v, want %v", string(gotValue), string(wantValue))
	}
}

func TestBadgerTxCommit(t *testing.T) {
	ddb := openInMemoryDB(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()
	txnReadBefore := ddb.NewTransaction(false)
	defer txnReadBefore.Discard()
	tx, key, wantValue := &BadgerTx{txn: txnWrite}, "testKey", []byte("testValue")

	// write key
	setKeyValue(t, txnWrite, key, wantValue)
	// ensure key is not visible in a separate transaction yet
	ensureKeyNotFound(t, txnReadBefore, key)

	err := tx.Commit()
	if err != nil {
		t.Fatalf("BadgerTx.Commit() error = %v, want nil", err)
	}

	txnReadAfter := ddb.NewTransaction(false)
	defer txnReadAfter.Discard()
	// ensure write committed and key is visible in a separate transaction now
	assertKeyValue(t, txnReadAfter, key, wantValue)
}

func TestBadgerTxRollback(t *testing.T) {
	ddb := openInMemoryDB(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()
	txnReadBefore := ddb.NewTransaction(false)
	defer txnReadBefore.Discard()
	tx, key, wantValue := &BadgerTx{txn: txnWrite}, "testKey", []byte("testValue")

	// write key
	setKeyValue(t, txnWrite, key, wantValue)
	// ensure key is not visible in a separate transaction
	ensureKeyNotFound(t, txnReadBefore, key)

	tx.Rollback()

	txnReadAfter := ddb.NewTransaction(false)
	defer txnReadAfter.Discard()
	// ensure write rolled back and key does not exist
	ensureKeyNotFound(t, txnReadAfter, key)
}
