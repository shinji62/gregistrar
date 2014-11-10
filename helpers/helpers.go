package helpers

import (
	"net"
	"regexp"
)

//Return list of Local IP
func LocalsIP(ipRegexp string) ([]string, error) {

	addrs, err := net.InterfaceAddrs()

	//Removing localIp
	r, _ := regexp.CompilePOSIX(ipRegexp)

	if err != nil {
		return nil, err

	}

	ipList := []string{}

	for _, addr := range addrs {

		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}
		ipv4 := ip.To4()
		if ipv4 != nil && !r.Match([]byte(ipv4.String())) {
			ipList = append(ipList, ipv4.String())
		}

	}

	if len(ipList) == 0 {
		return nil, nil
	}
	return ipList, nil

}
