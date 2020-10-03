package badgerdb

import (
	"net"
	"reflect"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/google/uuid"
)

func ipNetMustParse(t *testing.T, s string) *net.IPNet {
	_, ret, err := net.ParseCIDR(s)
	if err != nil {
		t.Fatalf("net.ParseCIDR() error = %v, want nil", err)
	}
	return ret
}

func TestCIDRMapToJSON(t *testing.T) {
	tests := []struct {
		name    string
		c       cidrMap
		want    string
		wantErr bool
	}{
		{
			name: "one cidr",
			c: cidrMap{
				uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"): ipNetMustParse(t, "192.168.0.0/24"),
			},
			want:    `{"2eba5e83-d0c3-46f0-bbeb-884e62e19b62":"192.168.0.0/24"}`,
			wantErr: false,
		},
		{
			name:    "empty",
			c:       cidrMap{},
			want:    `{}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.toJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("cidrMap.toJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("cidrMap.toJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestCIDRMapFromJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    cidrMap
		wantErr bool
	}{
		{
			name: "one cidr",
			args: args{
				data: []byte(`{"2eba5e83-d0c3-46f0-bbeb-884e62e19b62":"192.168.0.0/24"}`),
			},
			want: cidrMap{
				uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"): ipNetMustParse(t, "192.168.0.0/24"),
			},
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				data: []byte(`{}`),
			},
			want:    cidrMap{},
			wantErr: false,
		},
		{
			name: "bad json type",
			args: args{
				data: []byte(`42`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad key",
			args: args{
				data: []byte(`{42:"192.168.0.0/24"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "key invalid uuid",
			args: args{
				data: []byte(`{"non uuid":"192.168.0.0/24"}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "value invalid cidr",
			args: args{
				data: []byte(`{"2eba5e83-d0c3-46f0-bbeb-884e62e19b62":"non cidr"}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			got, err := cidrMapFromJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("cidrMapFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cidrMapFromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBadgerTxGetCIDRMap(t *testing.T) {
	tests := []struct {
		name    string
		want    cidrMap
		wantErr bool
	}{
		{
			name: "get",
			want: cidrMap{
				uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"): ipNetMustParse(t, "192.168.0.0/24"),
			},
			wantErr: false,
		},
		{
			name:    "default",
			want:    cidrMap{},
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
			got, err := tx.getCIDRMap()
			if (err != nil) != tt.wantErr {
				t.Errorf("badgerTx.getCIDRMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("badgerTx.getCIDRMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBadgerTxSaveCIDRMap(t *testing.T) {
	ddb := openInMemoryDB(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()
	tx := &badgerTx{txn: txnWrite}
	err := tx.saveCIDRMap(cidrMap{
		uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"): ipNetMustParse(t, "192.168.0.0/24"),
	})
	if err != nil {
		t.Fatalf("badgerTx.saveCIDRMap() error = %v, want nil", err)
	}
	cupaloy.New(cupaloy.SnapshotFileExtension(".json")).SnapshotT(t, dumpData(t, tx))
}

func TestBadgerTxGetSubnetCIDRs(t *testing.T) {
	tests := []struct {
		name    string
		want    []*net.IPNet
		wantErr bool
	}{
		{
			name:    "default",
			want:    []*net.IPNet{},
			wantErr: false,
		},
		{
			name:    "get",
			want:    []*net.IPNet{ipNetMustParse(t, "192.168.0.0/24")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt // pin!
		t.Run(tt.name, func(t *testing.T) {
			ddb := openInMemoryDBWithData(t)
			txnRead := ddb.NewTransaction(false)
			defer txnRead.Discard()
			tr := &badgerTx{
				txn: txnRead,
			}
			got, err := tr.GetSubnetCIDRs()
			if (err != nil) != tt.wantErr {
				t.Errorf("badgerTx.GetSubnetCIDRs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("badgerTx.GetSubnetCIDRs() = %v, want %v", got, tt.want)
			}
		})
	}
}
