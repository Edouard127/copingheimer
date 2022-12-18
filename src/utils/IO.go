package utils

import (
	"edouard127/copingheimer/src/intf"
	"fmt"
	"io"
	"net"
	"os"
)

func ReadBlacklist(cfg *intf.Arguments) ([][2]net.IP, error) {
	if file, err := os.Open(cfg.BlacklistFile); err != nil {
		return nil, err
	} else {
		defer file.Close()
		var blacklist [][2]net.IP
		for {
			// Check if the line is valid
			if ips, err := readMask(file); err != nil {
				if err == io.EOF {
					break
				} else {
					return nil, err
				}
			} else {
				blacklist = append(blacklist, ips)
			}
		}
		return blacklist, nil
	}
}

func readMask(r io.Reader) ([2]net.IP, error) {
	var (
		ip   [4]byte
		mask int
	)
	if _, err := fmt.Fscanf(r, "%d.%d.%d.%d/%d\n", &ip[0], &ip[1], &ip[2], &ip[3], &mask); err != nil {
		return [2]net.IP{}, err
	}
	var reversedMask int
	for i := 0; i < mask; i++ {
		reversedMask += 1 << i
	}
	return [2]net.IP{ip[:], getNextIP(net.IP(ip[:]).Mask(net.CIDRMask(mask, 32)), 0)}, nil // TODO: Incorrect
}
