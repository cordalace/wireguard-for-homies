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

// cidrToNet removes IP info from CIDR, e.g. "127.0.0.1/8" -> "127.0.0.0/8", "192.168.0.32/24" -> "192.168.0.0/24".
func cidrToNet(c *net.IPNet) *net.IPNet {
	return &net.IPNet{IP: c.IP.Mask(c.Mask), Mask: c.Mask}
}

func (i *ip) ListAddrCIDRs() ([]*net.IPNet, error) {
	// Should we use `netlink.FAMILY_V4 | netlink.FAMILY_V6` instead of `netlink.FAMILY_ALL`?
	// Is `netlink.FAMILY_MPLS` ok?
	addrs, err := netlink.AddrList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, err
	}

	ret := make([]*net.IPNet, 0, len(addrs))
	for _, addr := range addrs {
		ret = append(ret, cidrToNet(addr.IPNet))
	}

	return ret, nil
}
