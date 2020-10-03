package badgerdb

import (
	"reflect"
	"testing"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	badger "github.com/dgraph-io/badger/v2"
)

func TestNewBadgerDB(t *testing.T) {
	opts := badger.DefaultOptions("").WithInMemory(true)
	want := &badgerDB{opts: opts}
	if got := NewBadgerDB(opts); !reflect.DeepEqual(got, want) {
		t.Fatalf("NewBadgerDB() = %v, want %v", got, want)
	}
}

func TestBadgerDBInit(t *testing.T) {
	d := &badgerDB{
		db:   nil,
		opts: badger.DefaultOptions("").WithInMemory(true),
	}
	if err := d.Init(); err != nil {
		t.Fatalf("badgerDB.Init() error = %v, wantErr nil", err)
	}

	if d.db == nil {
		t.Fatalf("badgerDB.Init(), db is nil")
	}

	if d.db.IsClosed() {
		t.Fatalf("badgerDB.Init(), d.db.IsClosed() is true, want false")
	}
}

func TestBadgerDBClose(t *testing.T) {
	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		t.Fatalf("badger.Open() error = %v, wantErr nil", err)
	}
	type fields struct {
		db   *badger.DB
		opts badger.Options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		wantDB  *badger.DB
	}{
		{
			name: "db is nil",
			fields: fields{
				db:   nil,
				opts: opts,
			},
			wantErr: false,
			wantDB:  nil,
		},
		{
			name: "close badger db",
			fields: fields{
				db:   db,
				opts: opts,
			},
			wantErr: false,
			wantDB:  nil,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			d := &badgerDB{
				db:   tt.fields.db,
				opts: tt.fields.opts,
			}
			if err := d.Close(); (err != nil) != tt.wantErr {
				t.Errorf("badgerDB.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(d.db, tt.wantDB) {
				t.Errorf("badgerDB.Close() badgerDB.db = %v, wantDB %v", d.db, tt.wantDB)
			}
		})
	}
}

func TestBadgerDBBegin(t *testing.T) {
	opts := badger.DefaultOptions("").WithInMemory(true)
	ddb, err := badger.Open(opts)
	if err != nil {
		t.Fatalf("badger.Open() error = %v, wantErr nil", err)
	}
	txnRead := ddb.NewTransaction(false)
	defer txnRead.Discard()
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()

	type fields struct {
		db   *badger.DB
		opts badger.Options
	}
	type args struct {
		mode db.TxMode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    db.Tx
		wantErr bool
	}{
		{
			name: "db is nil",
			fields: fields{
				db:   nil,
				opts: opts,
			},
			args: args{
				mode: db.TxModeReadOnly,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "read only",
			fields: fields{
				db:   ddb,
				opts: opts,
			},
			args: args{
				mode: db.TxModeReadOnly,
			},
			want:    &badgerTx{txn: txnRead},
			wantErr: false,
		},
		{
			name: "read write",
			fields: fields{
				db:   ddb,
				opts: opts,
			},
			args: args{
				mode: db.TxModeReadWrite,
			},
			want:    &badgerTx{txn: txnWrite},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			d := &badgerDB{
				db:   tt.fields.db,
				opts: tt.fields.opts,
			}
			got, err := d.Begin(tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("badgerDB.Begin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("badgerDB.Begin() = %v, want %v", got, tt.want)
			}
		})
	}
}
