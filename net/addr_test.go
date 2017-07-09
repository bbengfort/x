package net_test

import (
	"fmt"
	"net"
	"strconv"

	. "github.com/bbengfort/x/net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Addr", func() {

	Describe("ExternalIP", func() {

		var (
			ip  string
			err error
		)

		It("should return an IP address without an error", func() {
			ip, err = ExternalIP()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(ip).ShouldNot(BeNil())
		})

		It("should return an IP address that is parseable", func() {
			ip, err = ExternalIP()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(net.ParseIP(ip)).ShouldNot(BeNil())
		})

		It("should return an IP address that is not the loopback", func() {
			ip, err = ExternalIP()
			Ω(err).ShouldNot(HaveOccurred())

			addr := net.ParseIP(ip)
			Ω(addr).ShouldNot(BeNil())
			Ω(addr.IsLoopback()).Should(BeFalse())
		})

	})

	Describe("ResolveAddr", func() {

		var (
			ip   string
			port string
			err  error
		)

		BeforeEach(func() {
			ip, err = ExternalIP()
			Ω(err).ShouldNot(HaveOccurred())

			port = strconv.Itoa(DefaultPort)
		})

		It("should resolve an empty address", func() {
			addr, err := ResolveAddr("")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(addr).Should(Equal(fmt.Sprintf("%s:%s", ip, port)))
		})

		It("should return an error if port is missing in an address", func() {
			_, err := ResolveAddr("192.168.1.1")
			Ω(err).Should(HaveOccurred())
		})

		It("should add the default port to an IP address with port 0", func() {
			addr, err := ResolveAddr("192.168.1.1:0")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(addr).Should(Equal(fmt.Sprintf("192.168.1.1:%s", port)))
		})

		It("should add the external IP address to a port", func() {
			addr, err := ResolveAddr(":5356")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(addr).Should(Equal(fmt.Sprintf("%s:5356", ip)))
		})

	})

	Describe("PublicIP", func() {

		It("should fetch the external IP address of the host", func() {
			ipaddr, err := PublicIP()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(ipaddr).ShouldNot(BeEmpty())
		})

	})

})
