package test

import (
	"edouard127/copingheimer/src/intf"
	"fmt"
	provider "github.com/Tnze/go-mc/bot"
	"testing"
	"time"
)

func TestPinger(t *testing.T) {
	ip := "2b2t.org"
	if data, _, err := provider.PingAndListTimeout(ip, 5*time.Second); err != nil {
		t.Error(err)
	} else {
		status := &intf.StatusResponse{}
		if err := status.Put(data); err != nil {
			t.Error(err)
		} else {
			fmt.Println(status)
		}
	}
}
