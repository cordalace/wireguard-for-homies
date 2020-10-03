package badgerdb

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func TestBadgerTxGetOrCreateDeviceName(t *testing.T) {
	type args struct {
		defaultDeviceName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "default",
			args:    args{defaultDeviceName: "defaultWg0"},
			want:    "defaultWg0",
			wantErr: false,
		},
		{
			name:    "get",
			args:    args{defaultDeviceName: "defaultWg0"},
			want:    "otherWg0",
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
			got, err := tx.GetOrCreateDeviceName(tt.args.defaultDeviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("badgerTx.GetOrCreateDeviceName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("badgerTx.GetOrCreateDeviceName() = %v, want %v", got, tt.want)
			}
			cupaloy.New(cupaloy.SnapshotFileExtension(".json")).SnapshotT(t, dumpData(t, tx))
		})
	}
}
