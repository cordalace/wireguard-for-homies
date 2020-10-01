package db

import (
	"errors"
	"net"

	"github.com/cordalace/wireguard-for-homies/internal/models"
	"github.com/google/uuid"
)

var ErrSubnetNotFound = errors.New("subnet not found")

type DB interface {
	Init() error
	Begin() (Tx, error)
	Close() error
}

type Tx interface {
	Commit() error
	Rollback()
	GetOrCreateDeviceName(defaultDeviceName string) (string, error)
	CreateSubnet(subnet *models.Subnet) (*models.Subnet, error)
	GetSubnet(id uuid.UUID) (*models.Subnet, error)
	DeleteSubnet(id uuid.UUID) error
	GetSubnetCIDRs() ([]*net.IPNet, error)
}
