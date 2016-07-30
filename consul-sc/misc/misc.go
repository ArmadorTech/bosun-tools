package misc

import (
	"net"
	"os"
)

var (
	myIP []net.IPAddr
)

func LocalIPs() []net.IP {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	ips := make([]net.IP,0)
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips=append(ips, ipnet.IP)
			}
		}
	}

	return ips
}

func LocalHostname() string {
	host, e := os.Hostname()
	if nil != e {
		return ""
	}
	return host
}
