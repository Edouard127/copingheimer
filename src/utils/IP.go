package utils

import (
	"edouard127/copingheimer/src/intf"
	"math"
	"math/rand"
	"net"
)

func getNextIP(ip net.IP, offset int) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += uint(offset)
	if v >= 167772160 && v <= 184549375 {
		v += 16777216
	} else if v >= 2886729728 && v <= 2887778303 {
		v += 1048576
	} else if v >= 3232235520 && v <= 3232301055 {
		v += 65536
	} else if v >= 2130706432 && v <= 2147483647 {
		v += 16777216
	} else if v >= 2851995648 && v <= 2852061183 {
		v += 1048576
	} else if v >= 2885681152 && v <= 2886746111 {
		v += 1048576
	}
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return []byte{v0, v1, v2, v3}
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
	return getNextIP(ip, 0)
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
	return getNextIP(s.CurIP, i)
}

func (s SubnetIterator) SetCurrent(ip net.IP) {
	s.CurIP = ip
}
