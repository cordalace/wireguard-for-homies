package badgerdb

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/cordalace/wireguard-for-homies/internal/models"
	"github.com/google/uuid"
)

func TestBadgerTxCreateSubnet(t *testing.T) {
	want := &models.Subnet{
		ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
		CIDR: ipNetMustParse(t, "192.168.0.0/24"),
	}
	ddb := openInMemoryDB(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()

	tx := &badgerTx{txn: txnWrite}
	got, err := tx.CreateSubnet(want)
	if err != nil {
		t.Fatalf("badgerTx.CreateSubnet() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("badgerTx.CreateSubnet() = %v, want %v", got, want)
	}

	cupaloy.New(cupaloy.SnapshotFileExtension(".json")).SnapshotT(t, dumpData(t, tx))
}

func TestBadgerTxCreateSubnetDuplicate(t *testing.T) {
	want := &models.Subnet{
		ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
		CIDR: ipNetMustParse(t, "192.168.0.0/24"),
	}
	ddb := openInMemoryDBWithData(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()

	tx := &badgerTx{txn: txnWrite}
	got, err := tx.CreateSubnet(want)
	if !errors.Is(err, db.ErrAlreadyExists) {
		t.Fatalf("badgerTx.CreateSubnet() error = %v, want db.ErrAlreadyExists", err)
	}
	if got != nil {
		t.Errorf("badgerTx.CreateSubnet() = %v, want nil", got)
	}
}

func TestBadgerTxGetSubnet(t *testing.T) {
	want := &models.Subnet{
		ID:   uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"),
		CIDR: ipNetMustParse(t, "192.168.0.0/24"),
	}
	ddb := openInMemoryDBWithData(t)
	txnRead := ddb.NewTransaction(false)
	defer txnRead.Discard()

	tx := &badgerTx{txn: txnRead}
	got, err := tx.GetSubnet(want.ID)
	if err != nil {
		t.Fatalf("badgerTx.GetSubnet() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("badgerTx.GetSubnet() = %v, want %v", got, want)
	}
}

func TestBadgerTxGetSubnetNotFound(t *testing.T) {
	ddb := openInMemoryDB(t)
	txnRead := ddb.NewTransaction(false)
	defer txnRead.Discard()

	tx := &badgerTx{txn: txnRead}
	got, err := tx.GetSubnet(uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"))
	if !errors.Is(err, db.ErrNotFound) {
		t.Fatalf("badgerTx.GetSubnet() error = %v, want db.ErrNotFound", err)
	}
	if got != nil {
		t.Errorf("badgerTx.GetSubnet() = %v, want nil", got)
	}
}

func TestBadgerDeleteSubnet(t *testing.T) {
	ddb := openInMemoryDBWithData(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()

	tx := &badgerTx{txn: txnWrite}
	err := tx.DeleteSubnet(uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"))
	if err != nil {
		t.Fatalf("badgerTx.DeleteSubnet() error = %v, want nil", err)
	}

	cupaloy.New(cupaloy.SnapshotFileExtension(".json")).SnapshotT(t, dumpData(t, tx))
}

func TestBadgerDeleteSubnetNotFound(t *testing.T) {
	ddb := openInMemoryDB(t)
	txnWrite := ddb.NewTransaction(true)
	defer txnWrite.Discard()

	tx := &badgerTx{txn: txnWrite}
	err := tx.DeleteSubnet(uuid.MustParse("2eba5e83-d0c3-46f0-bbeb-884e62e19b62"))
	if !errors.Is(err, db.ErrNotFound) {
		t.Fatalf("badgerTx.DeleteSubnet() error = %v, want db.ErrNotFound", err)
	}
}
