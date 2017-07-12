package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

// DefaultPort is used to allocate a port for services on the network when
// not specified directly. The port number refers to our office in A.V.
// Williams during graduate school and we've commonly used it in many
// applications. The DefaultPort is primarily used in ResolveAddr.
const DefaultPort = 3264

// ExternalIP looks up an the first available external IP address used by
// local network interfaces. This function returns a string representation of
// the IP address that can be parsed by net.IPAddr or other tools.
//
// NOTE: this function does not refer to the external network IP address, but
// rather non-local or loopback addresses on the machine. If the machine
// receives an internal DHCP provied IP address, then this function will
// detect that, not the IP address of the router. To find the publically
// accessible IP address of the machine use PublicIP.
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

// PublicIP makes an external HTTP request to myexternalip.com in order to
// discover the publically available IP address of the machine. This is
// especially useful when the machine sits behind a NAT device such as a
// router that performs port forwarding.
//
// NOTE: the myexternalip.com service maintains a rate limit of 30 requests
// per minute, do not exceed it!
func PublicIP() (string, error) {
	// Conduct the request with a 5 second timeout
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Get("http://ipv4.myexternalip.com/json")
	if err != nil {
		return "", err
	}

	// Ensure connection is closed on complete
	defer resp.Body.Close()

	// Check the status from the client
	if resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			return "", fmt.Errorf(
				"received staus %s: rate limit of 30 requests per minute exceeded",
				resp.Status,
			)
		}

		return "", fmt.Errorf(
			"could not lookup public IP address: %s", resp.Status,
		)
	}

	// Parse the body of the response
	data := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	// Get the IP address
	ipaddr, ok := data["ip"]
	if !ok {
		return "", errors.New("could not find IP address in response")
	}

	return ipaddr.(string), nil

}
