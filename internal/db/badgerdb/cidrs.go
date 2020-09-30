package badgerdb

import (
	"encoding/json"
	"net"

	"github.com/google/uuid"
)

const cidrMapKey = "cidrMap"

type cidrMap map[uuid.UUID]*net.IPNet

type cidrMapJSON map[string]string

func (c cidrMap) toJSON() ([]byte, error) {
	j := make(cidrMapJSON, len(c))
	for subnetID, subnetCIDR := range c {
		j[subnetID.String()] = subnetCIDR.String()
	}
	return json.Marshal(j)
}

func cidrMapFromJSON(data []byte) (cidrMap, error) {
	j := make(cidrMapJSON)
	err := json.Unmarshal(data, &j)
	if err != nil {
		return nil, err
	}
	ret := make(cidrMap, len(j))
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

func (t *badgerTX) getCIDRMap() (cidrMap, error) {
	value, err := getOrCreate(t.txn, cidrMapKey, []byte("{}"))
	if err != nil {
		return nil, err
	}

	return cidrMapFromJSON(value)
}

func (t *badgerTX) saveCIDRMap(c cidrMap) error {
	data, err := c.toJSON()
	if err != nil {
		return err
	}
	return t.txn.Set([]byte(cidrMapKey), data)
}

func (t *badgerTX) GetSubnetCIDRs() ([]*net.IPNet, error) {
	value, err := getOrCreate(t.txn, cidrMapKey, []byte("{}"))
	if err != nil {
		return nil, err
	}

	cm, err := cidrMapFromJSON(value)
	if err != nil {
		return nil, err
	}

	ret := make([]*net.IPNet, 0, len(cm))
	for _, subnetCIDR := range cm {
		ret = append(ret, subnetCIDR)
	}

	return ret, nil
}
