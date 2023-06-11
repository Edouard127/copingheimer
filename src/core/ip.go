package core

import (
	"net"
)

var DefaultIP = net.IPv4(1, 0, 0, 0).To4()

func GetIP(ip net.IP, o uint) net.IP {
	return IntIP(skipPrivate(IPInt(ip) + o))
}

// skipPrivate will skip addresses that are private.
// for example, 10.0.0.0/8, 172.16.0.0/12 and 192.168.0.0/16
func skipPrivate(ip uint) uint {

	return ip
}

func IPin(ip net.IP, a net.IPNet) bool {
	return IPInt(ip) >= IPInt(a.IP) && IPInt(ip) <= IPInt(LastIP(a))
}

func LastIP(a net.IPNet) net.IP {
	return net.IPv4(
		a.IP[0]|^a.Mask[0],
		a.IP[1]|^a.Mask[1],
		a.IP[2]|^a.Mask[2],
		a.IP[3]|^a.Mask[3],
	)
}

func IPInt(ip net.IP) uint {
	return uint(ip[0])<<24 | uint(ip[1])<<16 | uint(ip[2])<<8 | uint(ip[3])
}

func IntIP(i uint) net.IP {
	return net.IPv4(byte(i>>24), byte(i>>16), byte(i>>8), byte(i)).To4()
}
