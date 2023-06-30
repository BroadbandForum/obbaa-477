package osadapters

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

// receive the wan interface to be excluded of returned available ifaces
func AvailableIfaces(wan string) []string {

	var listIfaces []string
	const mgnt = "ens3"

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	for _, Iface := range interfaces {

		if Iface.Flags&net.FlagUp == 1 && Iface.Flags&net.FlagLoopback == 0 && Iface.Name != wan && Iface.Name != mgnt && len(Iface.Name) < 5 {
			listIfaces = append(listIfaces, Iface.Name)
		}
	}

	return listIfaces
}

// Get IPv4 address for a string interface name
func GetInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		return "", errors.New(fmt.Sprintf("interface %s don't have an ipv4 address\n", interfaceName))
	}
	return ipv4Addr.String(), nil
}

// return the MAC addr of a string
func GetMAC(iface string) ([]byte, string) {
	var ifaceMAC []byte
	var ifaceMACString string
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	for _, interf := range interfaces {

		if interf.Name == iface {
			ifaceMACString = interf.HardwareAddr.String()
			ifaceMAC = interf.HardwareAddr
		}
	}
	return ifaceMAC, ifaceMACString
}

// return the hostname of device
func GetHostname() (hostname string, err error) {
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatal(errors.New(fmt.Sprintf("Unable to get hostname!\n")))
		os.Exit(1)
	}

	return hostname, nil
}
