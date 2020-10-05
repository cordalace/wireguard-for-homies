package badgerdb

import (
	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/cordalace/wireguard-for-homies/internal/models"
	badger "github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
)

const subnetPrefix = "subnet"

func (t *BadgerTx) CreateSubnet(subnet *models.Subnet) (*models.Subnet, error) {
	key := fmtDBKey(subnetPrefix, subnet.ID.String())

	exists, err := t.exists(key)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, db.ErrAlreadyExists
	}

	subnetJSON, err := subnet.ToJSON()
	if err != nil {
		return nil, err
	}
	err = t.txn.Set(key, subnetJSON)
	if err != nil {
		return nil, err
	}

	c, err := t.getCIDRMap()
	if err != nil {
		return nil, err
	}

	c[subnet.ID] = subnet.CIDR

	err = t.saveCIDRMap(c)
	if err != nil {
		return nil, err
	}

	return &models.Subnet{ID: subnet.ID, CIDR: subnet.CIDR}, nil
}

func (t *BadgerTx) GetSubnet(id uuid.UUID) (*models.Subnet, error) {
	item, err := t.txn.Get(fmtDBKey(subnetPrefix, id.String()))
	switch err {
	case badger.ErrKeyNotFound:
		return nil, db.ErrNotFound
	case nil:
		var (
			value   []byte
			copyErr error
		)
		value, copyErr = item.ValueCopy(value)
		if copyErr != nil {
			return nil, copyErr
		}
		return models.SubnetFromJSON(value)
	default:
		return nil, err
	}
}

func (t *BadgerTx) DeleteSubnet(id uuid.UUID) error {
	key := fmtDBKey(subnetPrefix, id.String())

	exists, err := t.exists(key)
	if err != nil {
		return err
	}
	if !exists {
		return db.ErrNotFound
	}

	err = t.txn.Delete(key)
	if err != nil {
		return err
	}

	c, err := t.getCIDRMap()
	if err != nil {
		return err
	}

	delete(c, id)
	return t.saveCIDRMap(c)
}
