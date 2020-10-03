package badgerdb

import (
	"errors"
	"reflect"
	"testing"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/cordalace/wireguard-for-homies/internal/models"
	badger "github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
)

func TestBadgerTxCreateSubnet(t *testing.T) {
	withTestTx(t, initDBEmpty, txModeReadWrite, func(txn *badger.Txn) {
		want := &models.Subnet{
			ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
			CIDR: ipNetMustParse(t, "192.168.0.0/24"),
		}
		tx := &badgerTx{txn: txn}
		got, err := tx.CreateSubnet(want)
		if err != nil {
			t.Fatalf("badgerTx.CreateSubnet() error = %v, want nil", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("badgerTx.CreateSubnet() = %v, want %v", got, want)
		}
	})
}

func TestBadgerTxCreateSubnetDuplicate(t *testing.T) {
	withTestTx(t, initDBWithInput, txModeReadWrite, func(txn *badger.Txn) {
		want := &models.Subnet{
			ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
			CIDR: ipNetMustParse(t, "192.168.0.0/24"),
		}
		tx := &badgerTx{txn: txn}
		got, err := tx.CreateSubnet(want)
		if !errors.Is(err, db.ErrAlreadyExists) {
			t.Fatalf("badgerTx.CreateSubnet() error = %v, want db.ErrAlreadyExists", err)
		}
		if got != nil {
			t.Errorf("badgerTx.CreateSubnet() = %v, want nil", got)
		}
	})
}

func TestBadgerTxGetSubnet(t *testing.T) {
	withTestTx(t, initDBWithInput, txModeReadOnly, func(txn *badger.Txn) {
		want := &models.Subnet{
			ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
			CIDR: ipNetMustParse(t, "192.168.0.0/24"),
		}
		tx := &badgerTx{txn: txn}
		got, err := tx.GetSubnet(want.ID)
		if err != nil {
			t.Fatalf("badgerTx.GetSubnet() error = %v, want nil", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("badgerTx.GetSubnet() = %v, want %v", got, want)
		}
	})
}

func TestBadgerTxGetSubnetNotFound(t *testing.T) {
	withTestTx(t, initDBEmpty, txModeReadOnly, func(txn *badger.Txn) {
		tx := &badgerTx{txn: txn}
		got, err := tx.GetSubnet(uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"))
		if !errors.Is(err, db.ErrNotFound) {
			t.Fatalf("badgerTx.GetSubnet() error = %v, want db.ErrNotFound", err)
		}
		if got != nil {
			t.Errorf("badgerTx.GetSubnet() = %v, want nil", got)
		}
	})
}

func TestBadgerDeleteSubnet(t *testing.T) {
	withTestTx(t, initDBWithInput, txModeReadWrite, func(txn *badger.Txn) {
		tx := &badgerTx{txn: txn}
		err := tx.DeleteSubnet(uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"))
		if err != nil {
			t.Fatalf("badgerTx.DeleteSubnet() error = %v, want nil", err)
		}
	})
}

func TestBadgerDeleteSubnetNotFound(t *testing.T) {
	withTestTx(t, initDBEmpty, txModeReadWrite, func(txn *badger.Txn) {
		tx := &badgerTx{txn: txn}
		err := tx.DeleteSubnet(uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"))
		if !errors.Is(err, db.ErrNotFound) {
			t.Fatalf("badgerTx.DeleteSubnet() error = %v, want db.ErrNotFound", err)
		}
	})
}
