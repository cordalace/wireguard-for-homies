package badgerdb

import (
	"reflect"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func TestFmtDBKey(t *testing.T) {
	expected := "testPrefix:testID"
	actual := string(fmtDBKey("testPrefix", "testID"))
	if actual != expected {
		t.Errorf("fmtDBKey() = %v, want %v", actual, expected)
	}
}

func TestBadgerTxGetOrCreate(t *testing.T) {
	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "create",
			args: args{
				key:   "testKey",
				value: []byte(`"test value"`),
			},
			want:    []byte(`"test value"`),
			wantErr: false,
		},
		{
			name: "get",
			args: args{
				key:   "testKey",
				value: []byte(`"test value"`),
			},
			want:    []byte(`"other value"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			ddb := openInMemoryDBWithData(t)

			txnWrite := ddb.NewTransaction(true)
			defer txnWrite.Discard()
			tx := &badgerTx{
				txn: txnWrite,
			}
			got, err := tx.getOrCreate(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("badgerTx.getOrCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("badgerTx.getOrCreate() = %v, want %v", string(got), string(tt.want))
			}
			cupaloy.New(cupaloy.SnapshotFileExtension(".json")).SnapshotT(t, dumpData(t, tx))
		})
	}
}

func TestBadgerTxGetOrDefault(t *testing.T) {
	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				key:   "testKey",
				value: []byte(`"test value"`),
			},
			want:    []byte(`"test value"`),
			wantErr: false,
		},
		{
			name: "get",
			args: args{
				key:   "testKey",
				value: []byte(`"test value"`),
			},
			want:    []byte(`"other value"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			ddb := openInMemoryDBWithData(t)

			txnRead := ddb.NewTransaction(false)
			defer txnRead.Discard()
			tx := &badgerTx{
				txn: txnRead,
			}
			got, err := tx.getOrDefault(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("badgerTx.getOrCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("badgerTx.getOrCreate() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
