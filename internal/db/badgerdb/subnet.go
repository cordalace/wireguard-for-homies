package badgerdb

import (
	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/cordalace/wireguard-for-homies/internal/models"
	badger "github.com/dgraph-io/badger/v2"
	"github.com/google/uuid"
)

const subnetPrefix = "subnet"

func (t *badgerTX) CreateSubnet(subnet *models.Subnet) (*models.Subnet, error) {
	var (
		subnetID uuid.UUID
		err      error
	)
	if subnet.ID == uuid.Nil {
		subnetID, err = uuid.NewRandom()
		if err != nil {
			return nil, err
		}
	} else {
		subnetID = subnet.ID
	}
	subnetJSON, err := subnet.ToJSON()
	if err != nil {
		return nil, err
	}
	err = t.txn.Set(fmtDBKey(subnetPrefix, subnetID.String()), subnetJSON)
	if err != nil {
		return nil, err
	}

	c, err := t.getCIDRMap()
	if err != nil {
		return nil, err
	}

	c[subnetID] = subnet.CIDR

	err = t.saveCIDRMap(c)
	if err != nil {
		return nil, err
	}

	return &models.Subnet{ID: subnetID, CIDR: subnet.CIDR}, nil
}

func (t *badgerTX) GetSubnet(id uuid.UUID) (*models.Subnet, error) {
	item, err := t.txn.Get(fmtDBKey(subnetPrefix, id.String()))
	switch err {
	case badger.ErrKeyNotFound:
		return nil, db.ErrSubnetNotFound
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

func (t *badgerTX) DeleteSubnet(id uuid.UUID) error {
	err := t.txn.Delete(fmtDBKey(subnetPrefix, id.String()))
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
