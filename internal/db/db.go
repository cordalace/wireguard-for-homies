package db

import (
	"errors"
	"net"

	"github.com/cordalace/wireguard-for-homies/internal/models"
	"github.com/google/uuid"
)

var ErrSubnetNotFound = errors.New("subnet not found")

type TxMode int

const (
	TxModeReadOnly TxMode = iota
	TxModeReadWrite
)

type DB interface {
	Init() error
	Close() error
	Begin(mode TxMode) (Tx, error)
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
