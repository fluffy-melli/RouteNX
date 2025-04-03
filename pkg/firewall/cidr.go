package firewall

import (
	"net"

	"github.com/fluffy-melli/RouteNX/pkg/config"
)

func IsCidrBlock(c *config.RouteNX, route *config.Route, ip string) bool {
	for _, firewall := range c.Firewall {
		for _, RoutefirewallL := range route.Firewall {
			if firewall.Name == RoutefirewallL {
				for _, blockIP := range firewall.CIDR {
					_, cidr, _ := net.ParseCIDR(blockIP)
					ipr := net.ParseIP(ip)
					if firewall.Block && cidr.Contains(ipr) {
						return true
					}
					if !firewall.Block && cidr.Contains(ipr) {
						return false
					}
				}
				if !firewall.Block {
					return true
				}
			}
		}
	}
	return false
}
