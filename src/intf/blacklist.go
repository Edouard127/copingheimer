package intf

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Blacklist [][2]net.IP

func ReadBlacklist(cfg *Arguments) (*Blacklist, error) {
	b := &Blacklist{}
	if file, err := os.Open(cfg.BlacklistFile); err != nil {
		return nil, err
	} else {
		defer file.Close()
		for {
			// Check if the line is valid
			if ips, err := b.ReadFrom(file); err != nil {
				if err == io.EOF {
					break
				} else {
					return nil, err
				}
			} else {
				*b = append(*b, ips)
			}
		}
		return b, nil
	}
}

func (b *Blacklist) ReadFrom(r io.Reader) ([2]net.IP, error) {
	var (
		ip   [4]byte
		mask int
	)
	if _, err := fmt.Fscanf(r, "%d.%d.%d.%d/%d\n", &ip[0], &ip[1], &ip[2], &ip[3], &mask); err != nil {
		return [2]net.IP{}, err
	}
	ipAddr := net.IP(ip[:])
	maskAddr := net.CIDRMask(mask, 32)
	lastIP := net.IPv4(
		ipAddr[0]|^maskAddr[0],
		ipAddr[1]|^maskAddr[1],
		ipAddr[2]|^maskAddr[2],
		ipAddr[3]|^maskAddr[3],
	)
	return [2]net.IP{ip[:], lastIP}, nil
}
