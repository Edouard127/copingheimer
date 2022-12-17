package main

import (
	"edouard127/copingheimer/src/intf"
	"edouard127/copingheimer/src/utils"
	"fmt"
	provider "github.com/Tnze/go-mc/bot"
	"github.com/boltdb/bolt"
	"net"
	"path"
	"time"
)

var (
	Config = intf.LoadConfig("config.env")
	Subnet = &net.IPNet{
		IP:   net.IP{193, 0, 0, 0},
		Mask: net.IPMask{255, 255, 0, 0},
	}
	IP            = utils.IPSubnetIterator(Subnet)
	Database, err = bolt.Open(path.Join(Config.Database, "copingheimer.db"), 0600, nil)
	tasks         = 0
)

func main() {
	if err != nil {
		panic(err)
	}
	fmt.Println(
		"Copingheimer, by Kamigen\n" +
			"This program will ping Minecraft servers all around the world and try to get as much data as possible",
	)

	if err := Database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("servers"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		panic(err)
	}

	for {
		// If the CPUSaver is enabled, we will only start a new instance if there is one available
		if Config.CPUSaver {
			if tasks < Config.InstanceCount {
				tasks++
				go HandleServer(tasks)
			}
		} else {
			go HandleServer(tasks)
			tasks++
		}
	}
}

func HandleServer(instance int) {
	ip := IP().GetNext(instance)
	// Ping the server
	if data, _, err := provider.PingAndListTimeout(ip.String(), time.Duration(Config.Timeout)*time.Millisecond); err != nil {
		fmt.Println("failed to ping", ":", err)
	} else {
		status := &intf.StatusResponse{}
		if err := status.ReadFrom(data); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(status)
			// Write to database
			if err := Database.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("servers"))
				err := b.Put([]byte(IP().CurIP.String()), data)
				return err
			}); err != nil {
				panic(err)
			}
		}
	}
	tasks--
	IP().Next()
}
