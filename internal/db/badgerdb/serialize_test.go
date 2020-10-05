package badgerdb

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	badger "github.com/dgraph-io/badger/v2"
)

func dumpData(t *testing.T, tx *BadgerTx) []byte {
	data, err := tx.DumpData()
	if err != nil {
		t.Fatalf("BadgerTx.DumpData() error = %v, want nil", err)
	}
	return data
}

func TestBadgerTxDumpData(t *testing.T) {
	ddb := openInMemoryDB(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()
	setKeyValue(t, txnWrite, "testKey", []byte(`"test value"`))
	err := txnWrite.Commit()
	if err != nil {
		t.Errorf("badger.Txn.Commit() error = %v, want nil", err)
	}

	txnRead := ddb.NewTransaction(false)
	defer txnRead.Discard()
	tx := BadgerTx{txn: txnRead}
	data, err := tx.DumpData()
	if err != nil {
		t.Errorf("BadgerTx.DumpData() error = %v, want nil", err)
	}

	cupaloy.New(cupaloy.SnapshotFileExtension(".json")).SnapshotT(t, data)
}

func TestBadgerTxLoadDataDBState(t *testing.T) {
	ddb := openInMemoryDB(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()
	tx := BadgerTx{txn: txnWrite}
	err := tx.LoadData([]byte(`{"testKey":"test value"}`))
	if err != nil {
		t.Errorf("BadgerTx.LoadData() error = %v, want nil", err)
	}
	err = txnWrite.Commit()
	if err != nil {
		t.Errorf("badger.Txn.Commit() error = %v, want nil", err)
	}

	txnRead := ddb.NewTransaction(false)
	defer txnRead.Discard()
	assertKeyValue(t, txnRead, "testKey", []byte(`"test value"`))
}

func TestBadgerTxLoadData(t *testing.T) {
	ddb := openInMemoryDB(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()

	type fields struct {
		txn *badger.Txn
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{txn: txnWrite},
			args:    args{data: []byte(`{"testKey":"test value"}`)},
			wantErr: false,
		},
		{
			name:    "json syntax error",
			fields:  fields{txn: txnWrite},
			args:    args{data: []byte("invalid json")},
			wantErr: true,
		},
		{
			name:    "json unexpected type",
			fields:  fields{txn: txnWrite},
			args:    args{data: []byte(`"non json object type"`)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			tr := &BadgerTx{
				txn: tt.fields.txn,
			}
			if err := tr.LoadData(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("BadgerTx.LoadData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
