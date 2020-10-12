package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/google/uuid"
)

var ErrSubnetOverlaps = errors.New("subnet overlaps with one of existing subnets")

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

func subnetCIDRsEqual(a, b *net.IPNet) bool {
	return a.IP.Equal(b.IP) && bytes.Equal(a.Mask, b.Mask)
}

// getSystemSubnetCIDRs returns system networks/subnets. It works by collecting addresses of network devices via
// NETLINK and removing client subnets from that collection.
func (m *Manager) getSystemSubnetCIDRs(clientCIDRs []*net.IPNet) ([]*net.IPNet, error) {
	systemCIDRs, err := m.ip.ListAddrCIDRs()
	if err != nil {
		return nil, err
	}

	ret := make([]*net.IPNet, 0, len(systemCIDRs))
	for _, systemCIDR := range systemCIDRs {
		found := false
		for _, clientCIDR := range clientCIDRs {
			if subnetCIDRsEqual(systemCIDR, clientCIDR) {
				found = true
				break
			}
		}
		if !found {
			ret = append(ret, systemCIDR)
		}
	}

	return ret, nil
}

func (m *Manager) CreateSubnet(subnet *Subnet) (*Subnet, error) {
	tx, err := m.db.Begin(db.TxModeReadWrite)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	clientCIDRs, err := tx.GetSubnetCIDRs()
	if err != nil {
		return nil, err
	}

	// As there is no client isolation (all clients use the same wireguard network interface) we should not allow to
	// create overlapping subnets
	if err := cidr.VerifyNoOverlap(clientCIDRs, subnet.CIDR); err != nil {
		// given subnet overlaps with one of client subnets
		return nil, ErrSubnetOverlaps
	}

	// User should not be able create subnet that overlaps with one of system subnet, system subnets are subnets used
	// by Linux operating system itself (or by system administrator) for any technical reason, e.g.:
	// - loopback interface network 127.0.0.1/8
	// - docker's eth0@ifXX network 172.17.0.0/16 used for internet access
	// - other networks created by system administrators, automation scripts, etc
	systemCIDRs, err := m.getSystemSubnetCIDRs(clientCIDRs)
	if err != nil {
		return nil, err
	}
	if err := cidr.VerifyNoOverlap(systemCIDRs, subnet.CIDR); err != nil {
		// given subnet overlaps with one of system subnets
		return nil, ErrSubnetOverlaps
	}

	ret, err := tx.CreateSubnet(subnet)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (m *Manager) GetSystemSubnetCIDRs() ([]*net.IPNet, error) {
	tx, err := m.db.Begin(db.TxModeReadOnly)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	clientCIDRs, err := tx.GetSubnetCIDRs()
	if err != nil {
		return nil, err
	}

	return m.getSystemSubnetCIDRs(clientCIDRs)
}
