package badgerdb

import (
	"errors"
	"testing"

	badger "github.com/dgraph-io/badger/v2"
)

func openInMemoryDB(t *testing.T) *badger.DB {
	opts := badger.DefaultOptions("").WithInMemory(true)
	ddb, err := badger.Open(opts)
	if err != nil {
		t.Fatalf("badger.Open() error = %v, want nil", err)
	}
	return ddb
}

func setKeyValue(t *testing.T, txn *badger.Txn, key, value string) {
	err := txn.Set([]byte(key), []byte(value))
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

func assertKeyValue(t *testing.T, txn *badger.Txn, key, wantValue string) {
	item, err := txn.Get([]byte(key))
	if err != nil {
		t.Fatalf("badger.Txn.Get() error = %v, want nil", err)
	}

	var gotValueBytes []byte
	gotValueBytes, err = item.ValueCopy(gotValueBytes)
	if err != nil {
		t.Fatalf("badger.Item.ValueCopy() error = %v, want nil", err)
	}

	gotValue := string(gotValueBytes)
	if gotValue != wantValue {
		t.Fatalf("badger.Txn.Get() = %v, want %v", gotValue, wantValue)
	}
}

func TestBadgerTxCommit(t *testing.T) {
	ddb := openInMemoryDB(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()
	txnReadBefore := ddb.NewTransaction(false)
	defer txnReadBefore.Discard()
	tx, key, wantValue := &badgerTx{txn: txnWrite}, "testKey", "testValue"

	// write key
	setKeyValue(t, txnWrite, key, wantValue)
	// ensure key written
	assertKeyValue(t, txnWrite, key, wantValue)
	// ensure key is not visible in a separate transaction yet
	ensureKeyNotFound(t, txnReadBefore, key)

	err := tx.Commit()
	if err != nil {
		t.Fatalf("badgerTx.Commit() error = %v, want nil", err)
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
	tx, key, wantValue := &badgerTx{txn: txnWrite}, "testKey", "testValue"

	// write key
	setKeyValue(t, txnWrite, key, wantValue)
	// ensure key written
	assertKeyValue(t, txnWrite, key, wantValue)
	// ensure key is not visible in a separate transaction
	ensureKeyNotFound(t, txnReadBefore, key)

	tx.Rollback()

	txnReadAfter := ddb.NewTransaction(false)
	defer txnReadAfter.Discard()
	// ensure write rolled back and key does not exist
	ensureKeyNotFound(t, txnReadAfter, key)
}
