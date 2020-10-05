package ip

import (
	"net"

	"github.com/vishvananda/netlink"
)

type IP interface {
	Init() error
	Close()
	LinkAddWg(linkName string) error
	LinkDelWg(linkName string) error
	AddrAdd(linkName string, addr *net.IPNet) error
}

func NewIP() IP {
	return &ip{}
}

type ip struct {
	netlinkHandle *netlink.Handle
}

func (i *ip) Init() error {
	var err error
	i.netlinkHandle, err = netlink.NewHandle()
	return err
}

func (i *ip) Close() {
	if i.netlinkHandle != nil {
		i.netlinkHandle.Delete()
		i.netlinkHandle = nil
	}
}
