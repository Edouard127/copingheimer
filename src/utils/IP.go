package utils

import (
	"net"
)

func getNextIP(ip net.IP, offset int) net.IP {
	ip = ip.To4()
	for i := 3; i >= 0; i-- {
		ip[i] += byte(offset)
		if ip[i] >= byte(offset) {
			break
		}
		offset = 1
	}
	return ip
}

func IPSubnetIterator(subnet *net.IPNet) func() SubnetIterator {
	var ip net.IP
	return func() SubnetIterator {
		if ip == nil {
			ip = subnet.IP
		} else {
			ip = getNextIP(ip, 1)
		}
		if subnet.Contains(ip) {
			return SubnetIterator{ip}
		}
		return SubnetIterator{nil}
	}
}

type SubnetIterator struct {
	CurIP net.IP
}

func (s SubnetIterator) Next() {
	s.CurIP = getNextIP(s.CurIP, 1)
}

func (s SubnetIterator) GetNext(i int) net.IP {
	return getNextIP(s.CurIP, i)
}

func (s SubnetIterator) NextSubnet(n int) {
	s.CurIP = getNextIP(s.CurIP, 256*n)
}
