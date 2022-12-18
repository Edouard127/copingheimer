package test

import (
	"edouard127/copingheimer/src/intf"
	"edouard127/copingheimer/src/utils"
	"net"
	"os"
	"path"
	"testing"
)

func Test_GetIn(t *testing.T) {
	ip, a, b := net.IP{192, 168, 0, 3}, net.IP{192, 168, 0, 1}, net.IP{192, 168, 0, 5}
	if utils.IPin(ip, a, b) {
		t.Log("IPin")
	} else {
		t.Log("Not IPin")
	}
}

func Test_Blacklist(t *testing.T) {
	pwd, _ := os.Getwd()
	if data, err := utils.ReadBlacklist(&intf.Arguments{
		BlacklistFile: path.Join(pwd, "test_blacklist.txt"),
	}); err != nil {
		t.Error(err)
	} else {
		for _, ip := range data {
			t.Log(ip)
		}
	}
}
