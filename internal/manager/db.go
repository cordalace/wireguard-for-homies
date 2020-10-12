package manager

import (
	"net"

	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/google/uuid"
)

type DB interface {
	db.DB
	Begin(mode db.TxMode) (Tx, error)
}

type Tx interface {
	db.Tx
	GetOrCreateDeviceName(defaultDeviceName string) (string, error)
	CreateSubnet(subnet *Subnet) (*Subnet, error)
	GetSubnet(id uuid.UUID) (*Subnet, error)
	DeleteSubnet(id uuid.UUID) error
	GetSubnetCIDRs() ([]*net.IPNet, error)
}
