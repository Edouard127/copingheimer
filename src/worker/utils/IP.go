package utils

import (
	"edouard127/copingheimer/src/intf"
	"github.com/valyala/fastrand"
	"math"
	"net"
)

func getNextIP(ip net.IP, offset int32, random bool) net.IP {
	i := ip.To4()
	if random {
		i[0] = byte(fastrand.Uint32n(255))
		i[1] = byte(fastrand.Uint32n(255))
		i[2] = byte(fastrand.Uint32n(255))
		i[3] = byte(fastrand.Uint32n(255))
		return i
	}
	v := IPInt(i) + offset
	return net.IPv4(byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func IPin(ip, a, b net.IP) bool {
	ip1 := segment(ip[0], 3) + segment(ip[1], 2) + segment(ip[2], 1) + segment(ip[3], 0)
	a1 := segment(a[0], 3) + segment(a[1], 2) + segment(a[2], 1) + segment(a[3], 0)
	b1 := segment(b[0], 3) + segment(b[1], 2) + segment(b[2], 1) + segment(b[3], 0)
	return ip1 >= a1 && ip1 <= b1
}

func IPInt(ip net.IP) int32 {
	return int32(ip[0])<<24 + int32(ip[1])<<16 + int32(ip[2])<<8 + int32(ip[3])
}

func segment(seg byte, exp uint) int {
	return int(seg) * int(math.Pow(256, float64(exp)))
}

func AddIP(a, b net.IP) net.IP {
	c := IPInt(a) + IPInt(b)
	return net.IPv4(byte(c>>24), byte(c>>16), byte(c>>8), byte(c))
}

func RandIP() net.IP {
	return getNextIP([]byte{0, 0, 0, 0}, 0, true)
}

func IPSubnetIterator(subnet net.IPNet, blacklist *intf.Blacklist) IteratorFunc {
	var ip net.IP
	return func(o int32, r bool) *SubnetIterator {
		if ip == nil {
			ip = subnet.IP
			ip = getNextIP(ip, o, r)
		} else {
			ip = getNextIP(ip, o, r)
		}
		if subnet.Contains(ip) {
			return &SubnetIterator{CurIP: ip, Blacklist: blacklist}
		}
		return &SubnetIterator{nil, blacklist}
	}
}

type IteratorFunc func(int32, bool) *SubnetIterator

type SubnetIterator struct {
	CurIP     net.IP
	Blacklist *intf.Blacklist
}

func (s SubnetIterator) Int() int32 {
	return IPInt(s.CurIP.To4())
}

func (s SubnetIterator) Array() []byte {
	return s.CurIP
}

func (s SubnetIterator) String() string {
	return s.CurIP.To4().String()
}
