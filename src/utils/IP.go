package utils

import (
	"edouard127/copingheimer/src/intf"
	"math"
	"math/rand"
	"net"
)

func getNextIP(ip net.IP, offset int) net.IP {
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	for i := 3; i >= 0; i-- {
		ip[i] += byte(offset)
		if ip[i] >= byte(offset) {
			break
		}
		offset = 1
	}
	return skipPrivate(ip)
}

func skipPrivate(ip net.IP) net.IP {
	if ip[0] == 10 {
		ip[0]++
	}
	if ip[0] == 172 && ip[1]&0xf0 == 16 {
		ip[1] = 32
	}
	if ip[0] == 192 && ip[1] == 168 {
		ip[0]++
		ip[1]++
	}
	return skipLoopback(ip)
}

func skipLoopback(ip net.IP) net.IP {
	if ip[0] == 127 {
		ip[0]++
	}
	return ip
}

func IPin(ip, a, b net.IP) bool {
	ip1 := segment(ip[0], 3) + segment(ip[1], 2) + segment(ip[2], 1) + segment(ip[3], 0)
	a1 := segment(a[0], 3) + segment(a[1], 2) + segment(a[2], 1) + segment(a[3], 0)
	b1 := segment(b[0], 3) + segment(b[1], 2) + segment(b[2], 1) + segment(b[3], 0)
	return ip1 >= a1 && ip1 <= b1
}

func segment(seg byte, exp uint) int {
	return int(seg) * int(math.Pow(256, float64(exp)))
}

func AddIP(a, b net.IP) net.IP {
	a1 := segment(a[0], 3) + segment(a[1], 2) + segment(a[2], 1) + segment(a[3], 0)
	b1 := segment(b[0], 3) + segment(b[1], 2) + segment(b[2], 1) + segment(b[3], 0)
	ipInt := a1 + b1
	p1 := ipInt / int(math.Pow(256, 3))
	p2 := ipInt % int(math.Pow(256, 3)) / int(math.Pow(256, 2))
	p3 := ipInt % int(math.Pow(256, 3)) / int(math.Pow(256, 2)) / int(math.Pow(256, 1))
	p4 := ipInt % int(math.Pow(256, 3)) / int(math.Pow(256, 2)) % int(math.Pow(256, 1))
	return net.IPv4(byte(p1), byte(p2), byte(p3), byte(p4))
}

func RandIP() net.IP {
	ip := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		ip[i] = byte(rand.Intn(255))
	}
	return skipPrivate(ip)
}

func IPSubnetIterator(subnet *net.IPNet, blacklist intf.Blacklist) func() SubnetIterator {
	var ip net.IP
	return func() SubnetIterator {
		if ip == nil {
			ip = subnet.IP
		} else {
			ip = getNextIP(ip, 1)
		}
		if subnet.Contains(ip) {
			return SubnetIterator{CurIP: ip, Blacklist: blacklist}
		}
		return SubnetIterator{nil, blacklist}
	}
}

type SubnetIterator struct {
	CurIP     net.IP
	Blacklist [][2]net.IP
}

func (s SubnetIterator) Next() {
	s.CurIP = getNextIP(s.CurIP, 1)
}

func (s SubnetIterator) GetNext(i int) net.IP {
	ip := getNextIP(s.CurIP, i)
	if s.Blacklist != nil {
		for _, ipRange := range s.Blacklist {
			if IPin(ip, ipRange[0], ipRange[1]) {
				return s.GetNext(i + 1)
			}
		}
	}
	return ip
}

func (s SubnetIterator) SetCurrent(ip net.IP) {
	s.CurIP = ip
}
