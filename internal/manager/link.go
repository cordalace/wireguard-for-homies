package manager

import "github.com/vishvananda/netlink"

// Wireguard represent links of type "wireguard", see https://www.wireguard.com/.
// It is backported from upstream, as no wireguard link present in v1.1.0, see: https://git.io/JUMTF
type wgLink struct {
	netlink.LinkAttrs
}

func (wg *wgLink) Attrs() *netlink.LinkAttrs {
	return &wg.LinkAttrs
}

func (wg *wgLink) Type() string {
	return "wireguard"
}
