package manager

import (
	"errors"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/cordalace/wireguard-for-homies/internal/models"
)

var ErrSubnetOverlaps = errors.New("subnet overlaps with one of existing subnets")

func (m *Manager) CreateSubnet(subnet *models.Subnet) (*models.Subnet, error) {
	tx, err := m.db.Begin(db.TxModeReadWrite)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	cidrs, err := tx.GetSubnetCIDRs()
	if err != nil {
		return nil, err
	}

	if err := cidr.VerifyNoOverlap(cidrs, subnet.CIDR); err != nil {
		return nil, ErrSubnetOverlaps
	}

	return tx.CreateSubnet(subnet)
}
