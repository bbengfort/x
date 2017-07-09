package net

import (
	"errors"
	"fmt"
	"net"
)

// DefaultPort is used to compute the TCP address in the absense of one.
const DefaultPort = 3265

// ExternalIP looks up an the first available external IP address used by
// local network interfaces. This function returns a string representation of
// the IP address that can be parsed by net.IPAddr or other tools.
//
// NOTE: this function does not refer to the external network IP address, but
// rather non-local or loopback addresses on the machine. If the machine
// receives an internal DHCP provied IP address, then this function will
// detect that, not the IP address of the router. To find the publically
// accessible IP address of the machine use PublicIP().
func ExternalIP() (string, error) {

	// Get addresses for the interface
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("Could not get interface addresses: %s", err.Error())
	}

	// Go through each address to find a an IPv4
	for _, addr := range addrs {

		var ip net.IP

		switch val := addr.(type) {
		case *net.IPNet:
			ip = val.IP
		case *net.IPAddr:
			ip = val.IP
		}

		if ip == nil || ip.IsLoopback() {
			continue // ignore loopback and nil addresses
		}

		ip = ip.To4()
		if ip == nil {
			continue // not an ipv4 address
		}

		return ip.String(), nil
	}

	return "", errors.New("Are you connected to the network?!")
}

// ResolveAddr accepts an address as a string and if the IP address is missing
// it replaces it with the result from ExternalIP then returns the addr
// string. Likewise if the Port is missing, it returns an address with the
// DefaultPort appended to the address string.
func ResolveAddr(addr string) (string, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("Could not resolve address: %s", err.Error())
	}

	if tcpAddr.IP == nil {
		ipstr, err := ExternalIP()
		if err != nil {
			return "", err
		}

		tcpAddr.IP = net.ParseIP(ipstr)
	}

	if tcpAddr.Port == 0 {
		tcpAddr.Port = DefaultPort
	}

	return tcpAddr.String(), nil
}
