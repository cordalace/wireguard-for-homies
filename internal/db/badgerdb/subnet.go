package badgerdb

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/cordalace/wireguard-for-homies/internal/models"
	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
)

type cidrs map[uuid.UUID]*net.IPNet

type cidrsJSON map[string]string

func (c cidrs) toJSON() ([]byte, error) {
	j := make(cidrsJSON, len(c))
	for subnetID, subnetCIDR := range c {
		j[subnetID.String()] = subnetCIDR.String()
	}
	return json.Marshal(j)
}

func cidrsFromJSON(data []byte) (cidrs, error) {
	j := make(cidrsJSON)
	err := json.Unmarshal(data, &j)
	if err != nil {
		return nil, err
	}
	ret := make(cidrs, len(j))
	var (
		parsedID   uuid.UUID
		parsedCIDR *net.IPNet
	)
	for key, value := range j {
		parsedID, err = uuid.Parse(key)
		if err != nil {
			return nil, err
		}
		_, parsedCIDR, err = net.ParseCIDR(value)
		if err != nil {
			return nil, err
		}
		ret[parsedID] = parsedCIDR
	}
	return ret, nil
}

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
	err = t.txn.Set([]byte(fmt.Sprintf("subnet:%v", subnetID.String())), subnetJSON)
	if err != nil {
		return nil, err
	}

	c, err := t.getCIDRs()
	if err != nil {
		return nil, err
	}

	c[subnetID] = subnet.CIDR

	err = t.saveCIDRs(c)
	if err != nil {
		return nil, err
	}

	return &models.Subnet{ID: subnetID, CIDR: subnet.CIDR}, nil
}

func (t *badgerTX) GetSubnet(id uuid.UUID) (*models.Subnet, error) {
	item, err := t.txn.Get([]byte(fmt.Sprintf("subnet:%v", id.String())))
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
	err := t.txn.Delete([]byte(id.String()))
	if err != nil {
		return err
	}

	c, err := t.getCIDRs()
	if err != nil {
		return err
	}

	delete(c, id)
	return t.saveCIDRs(c)
}

func (t *badgerTX) getCIDRs() (cidrs, error) {
	value, err := getOrCreate(t.txn, "cidrs", []byte("{}"))
	if err != nil {
		return nil, err
	}

	return cidrsFromJSON(value)
}

func (t *badgerTX) saveCIDRs(c cidrs) error {
	data, err := c.toJSON()
	if err != nil {
		return err
	}
	return t.txn.Set([]byte("cidrs"), data)
}

func (t *badgerTX) GetSubnetCIDRs() ([]*net.IPNet, error) {
	value, err := getOrCreate(t.txn, "cidrs", []byte("{}"))
	if err != nil {
		return nil, err
	}

	cidrs, err := cidrsFromJSON(value)
	if err != nil {
		return nil, err
	}

	ret := make([]*net.IPNet, 0, len(cidrs))
	for _, subnetCIDR := range cidrs {
		ret = append(ret, subnetCIDR)
	}

	return ret, nil
}
