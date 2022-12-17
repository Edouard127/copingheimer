package test

import (
	"edouard127/copingheimer/src/intf"
	"fmt"
	provider "github.com/Tnze/go-mc/bot"
	"github.com/boltdb/bolt"
	"testing"
	"time"
)

func TestPinger(t *testing.T) {
	Database, err := bolt.Open("copingheimer.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	ip := "2b2t.org"
	if err := Database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("servers"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		panic(err)
	}
	if data, _, err := provider.PingAndListTimeout(ip, 5*time.Second); err != nil {
		t.Error(err)
	} else {
		status := &intf.StatusResponse{}
		if err := status.ReadFrom(data); err != nil {
			t.Error(err)
		} else {
			t.Log(status)
			if err := Database.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("servers"))
				err := b.Put([]byte(ip), data)
				return err
			}); err != nil {
				panic(err)
			}
		}
	}
}
