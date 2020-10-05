package ip

import (
	"net"

	"github.com/vishvananda/netlink"
)

func (i *ip) AddrAdd(linkName string, addr *net.IPNet) error {
	a := &netlink.Addr{
		IPNet: addr,
	}
	return i.netlinkHandle.AddrAdd(getWgLink(linkName), a)
}
