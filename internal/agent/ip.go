package agent

import (
	"net"
)

func GetIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var l []net.IP
	var p []net.IP

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				l = append(l, ipnet.IP)
				continue
			}

			if ipnet.IP.IsPrivate() && ipnet.IP.To4() != nil {
				p = append(p, ipnet.IP)
				continue
			}

			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}

	if len(p) > 0 {
		return p[0], nil
	}

	if len(l) > 0 {
		return l[0], nil
	}

	return nil, nil
}
