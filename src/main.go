package main

import (
	"edouard127/copingheimer/src/intf"
	"edouard127/copingheimer/src/utils"
	"flag"
	"fmt"
	provider "github.com/Tnze/go-mc/bot"
	"net"
	"os"
	"time"
)

var (
	Config        = intf.LoadConfig("config.env")
	Arguments     = intf.Arguments{}
	Subnet        = &net.IPNet{}
	IP            = func() utils.SubnetIterator { return utils.SubnetIterator{} }
	Database, err = utils.InitDatabase(Config)
	tasks         = 0

	// Arguments
	// -c, --config: Path to the config file
	// -ip, --ip: IP address to start from
	// -i, --instances: Number of instances to run
	// -t, --timeout: Timeout for each ping
)

func init() {
	flag.StringVar(&Arguments.Config, "c", "config.env", "Path to the config file")
	flag.StringVar(&Arguments.IP, "ip", "0.0.0.1", "IP address to start from")
	flag.StringVar(&Arguments.IP, "i", "256", "Number of instances to run")
	flag.IntVar(&Arguments.Timeout, "t", 5000, "Timeout for each ping")
	flag.BoolVar(&Arguments.Help, "h", false, "Show this help")
	flag.Parse()

	if Arguments.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	Subnet = &net.IPNet{
		IP:   net.ParseIP(Arguments.IP),
		Mask: net.IPMask{0, 0, 0, 0},
	}
	IP = utils.IPSubnetIterator(Subnet)
}

func main() {
	if err != nil {
		panic(err)
	}
	fmt.Println(
		"Copingheimer, by Kamigen\n" +
			"This program will ping Minecraft servers all around the world and try to get as much data as possible",
	)

	// Read the current IP from the database
	if data, err := utils.Find(Database, "current_ip"); err != nil {
		fmt.Printf("failed to read current IP: %s\n", err)
	} else {
		if len(data) > 0 {
			IP().SetCurrent(net.ParseIP(data[0]))
		}
	}

	// Save the current IP on close
	defer func() {
		if err := utils.Write(Database, "current_ip", []byte(IP().CurIP.String())); err != nil {
			fmt.Println(err)
		}
	}()

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
			if err := utils.Write(Database, ip.String(), data); err != nil {
				fmt.Println(err)
			}
		}
	}
	tasks--
	IP().Next()
}
