package manager

import (
	"github.com/cordalace/wireguard-for-homies/internal/db"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

// Manager operates kernel wireguard settings.
type Manager struct {
	db            db.DB
	netlinkHandle *netlink.Handle
	logger        *zap.Logger
}

// NewWireguard creates new Wireguard instance.
func NewManager(db db.DB, netlinkHandle *netlink.Handle, logger *zap.Logger) *Manager {
	return &Manager{db: db, netlinkHandle: netlinkHandle, logger: logger}
}

func (m *Manager) getLink(name string) netlink.Link {
	la := netlink.NewLinkAttrs()
	la.Name = name
	wg := &wgLink{LinkAttrs: la}

	return wg
}

// Init wireguard manager, should be always called after instantiating Wireguard.
func (m *Manager) Init() error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	deviceName, err := tx.GetOrCreateDeviceName("wg0")
	if err != nil {
		return err
	}
	m.logger.Info("device name", zap.String("deviceName", deviceName))

	return netlink.LinkAdd(m.getLink(deviceName))
}

// Close wireguard.
func (m *Manager) Close() error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	deviceName, err := tx.GetOrCreateDeviceName("wg0")
	if err != nil {
		return err
	}
	m.logger.Info("device name", zap.String("deviceName", deviceName))

	return m.netlinkHandle.LinkDel(m.getLink(deviceName))
}
