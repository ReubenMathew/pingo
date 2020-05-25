package validate

import "net"

// Ipv4 : Validate a ipv4 address
func Ipv4(addr string) bool {
	return net.ParseIP(addr) != nil
}
