package ip

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

func getWgLink(name string) netlink.Link {
	la := netlink.NewLinkAttrs()
	la.Name = name
	wg := &wgLink{LinkAttrs: la}

	return wg
}

func (i *ip) LinkAddWg(linkName string) error {
	return i.netlinkHandle.LinkAdd(getWgLink(linkName))
}

func (i *ip) LinkDelWg(linkName string) error {
	return i.netlinkHandle.LinkDel(getWgLink(linkName))
}
