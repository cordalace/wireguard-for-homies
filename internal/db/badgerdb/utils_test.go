package badgerdb

import (
	"testing"

	badger "github.com/dgraph-io/badger/v2"
)

func TestFmtDBKey(t *testing.T) {
	expected := "testPrefix:testID"
	actual := string(fmtDBKey("testPrefix", "testID"))
	if actual != expected {
		t.Errorf("fmtDBKey() = %v, want %v", actual, expected)
	}
}

func TestGetOrCreateCreate(t *testing.T) {
	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		t.Fatalf("error creating in memory database: %v", err)
	}
	defer db.Close()

	txn := db.NewTransaction(true)
	defer txn.Discard()

	expected := "test value"
	actualBytes, err := getOrCreate(txn, "testKey", []byte(expected))
	if err != nil {
		t.Fatalf("error calling getOrCreate: %v", err)
	}

	actual := string(actualBytes)

	if actual != expected {
		t.Errorf("%v != %v", actual, expected)
	}
}
