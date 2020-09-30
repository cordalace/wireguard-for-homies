package models

import (
	"encoding/json"
	"net"

	"github.com/google/uuid"
)

type Subnet struct {
	ID   uuid.UUID
	CIDR *net.IPNet
}

type SubnetSlice []*Subnet

type subnetJSON struct {
	ID   string `json:"id"`
	CIDR string `json:"cidr"`
}

func (s *Subnet) ToJSON() ([]byte, error) {
	j := &subnetJSON{ID: s.ID.String(), CIDR: s.CIDR.String()}
	return json.Marshal(j)
}

func SubnetFromJSON(data []byte) (*Subnet, error) {
	j := &subnetJSON{}
	err := json.Unmarshal(data, j)
	if err != nil {
		return nil, err
	}
	ret := &Subnet{}
	ret.ID, err = uuid.Parse(j.ID)
	if err != nil {
		return nil, err
	}
	_, ret.CIDR, err = net.ParseCIDR(j.CIDR)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (l SubnetSlice) ToJSON() ([]byte, error) {
	j := make([]*subnetJSON, len(l))
	for i, item := range l {
		j[i] = &subnetJSON{ID: item.ID.String(), CIDR: item.CIDR.String()}
	}
	return json.Marshal(j)
}

func SubnetSliceFromJSON(data []byte) ([]*Subnet, error) {
	j := []*subnetJSON{}
	err := json.Unmarshal(data, &j)
	if err != nil {
		return nil, err
	}
	ret := make([]*Subnet, len(j))
	var (
		parsedID   uuid.UUID
		parsedCIDR *net.IPNet
	)
	for i, item := range j {
		parsedID, err = uuid.Parse(item.ID)
		if err != nil {
			return nil, err
		}
		_, parsedCIDR, err = net.ParseCIDR(item.CIDR)
		if err != nil {
			return nil, err
		}
		ret[i] = &Subnet{ID: parsedID, CIDR: parsedCIDR}
	}
	return ret, nil
}
