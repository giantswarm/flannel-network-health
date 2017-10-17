package key

import (
	"github.com/vishvananda/netlink"
)

func GetInterfaceIP(ipList []netlink.Addr) (string) {
	return ipList[0].IP.String()
}
