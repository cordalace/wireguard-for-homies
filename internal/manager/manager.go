package manager

import (
	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/cordalace/wireguard-for-homies/internal/ip"
	"go.uber.org/zap"
)

const defaultDeviceName = "wg0"

// Manager operates kernel wireguard settings.
type Manager struct {
	db     DB
	ip     ip.IP
	logger *zap.Logger
}

// NewWireguard creates new Wireguard instance.
func NewManager(db DB, i ip.IP, logger *zap.Logger) *Manager {
	return &Manager{db: db, ip: i, logger: logger}
}

// Init wireguard manager, should be always called after instantiating Wireguard.
func (m *Manager) Init() error {
	tx, err := m.db.Begin(db.TxModeReadWrite)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deviceName, err := tx.GetOrCreateDeviceName(defaultDeviceName)
	if err != nil {
		return err
	}

	return m.ip.LinkAddWg(deviceName)
}

// Close wireguard.
func (m *Manager) Close() error {
	tx, err := m.db.Begin(db.TxModeReadWrite)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deviceName, err := tx.GetOrCreateDeviceName(defaultDeviceName)
	if err != nil {
		return err
	}

	return m.ip.LinkDelWg(deviceName)
}
