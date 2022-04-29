package uniqueid

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

// createNodeID will return a number range from 0 to 1023
// It can be temp used for identifying a node with an unique id
func createNodeID() (int64, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		// Cannot run this application because the IP addresses cannot be found
		return 0, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				myIP = ipnet.IP.String()
				ipStr := strings.ReplaceAll(myIP, ".", "")
				ipInt, _ := strconv.ParseInt(ipStr, 10, 64)
				nodeID := ipInt % 1024
				return nodeID, nil
			}
		}
	}

	return 0, errors.New("error")
}
