package net

import "net"

// FreePort asks the kernel for a free, open port that is ready to use.
// https://github.com/phayes/freeport
func FreePort() (int, error) {

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer listen.Close()
	return listen.Addr().(*net.TCPAddr).Port, nil
}
