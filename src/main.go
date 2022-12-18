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
	flag.BoolVar(&Arguments.Help, "h", false, "Show this help message")
	flag.BoolVar(&Arguments.Help, "help", false, "Show this help message")
	flag.StringVar(&Arguments.Config, "c", "config.env", "Path to the config file")
	flag.StringVar(&Arguments.Config, "config", "config.env", "Path to the config file")
	flag.StringVar(&Arguments.Mode, "m", "random", "Mode to run in (default: \"random\") (random, order)")
	flag.StringVar(&Arguments.Mode, "mode", "random", "Mode to run in (default: \"random\") (random, order)")
	flag.StringVar(&Arguments.IP, "ip", "0.0.0.0", "IP address to start from with mask, only used in order mode")
	flag.BoolVar(&Arguments.CPUSaver, "cs", true, "Whether to enable the CPU saver or not (default: true)")
	flag.BoolVar(&Arguments.CPUSaver, "cpu-saver", true, "Whether to enable the CPU saver or not (default: true)")
	flag.IntVar(&Arguments.Instances, "i", 1, "Number of instances to run (default: 1)")
	flag.IntVar(&Arguments.Instances, "instances", 1, "Number of instances to run (default: 1)")
	flag.IntVar(&Arguments.Timeout, "t", 1000, "Timeout for each ping (default: 1000)")
	flag.IntVar(&Arguments.Timeout, "timeout", 1000, "Timeout for each ping (default: 1000)")
	flag.StringVar(&Arguments.BlacklistFile, "bf", "blacklist.txt", "Path to the blacklist file")
	flag.StringVar(&Arguments.BlacklistFile, "blacklist-file", "blacklist.txt", "Path to the blacklist file")
	flag.StringVar(&Arguments.Database, "d", "mongodb", "Database to use (default: mongodb) (mongodb, bolt)")
	flag.StringVar(&Arguments.Database, "database", "mongodb", "Database to use (default: mongodb) (mongodb, bolt)")
	flag.StringVar(&Arguments.DatabaseURL, "du", "mongodb://localhost:27017", "URL to the database")
	flag.StringVar(&Arguments.DatabaseURL, "database-url", "mongodb://localhost:27017", "URL to the database")
	flag.Parse()

	if Arguments.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if Config.InstanceCount == 0 {
		Config.InstanceCount = Arguments.Instances
	}
	if Config.Timeout == 0 {
		Config.Timeout = Arguments.Timeout
	}

	Subnet = &net.IPNet{
		IP:   net.ParseIP(Arguments.IP),
		Mask: net.IPMask{0, 0, 0, 0},
	}
	if blacklist, err := intf.ReadBlacklist(&Arguments); err == nil {
		IP = utils.IPSubnetIterator(Subnet, *blacklist)
	} else {
		fmt.Println("failed to read blacklist:", err)
		os.Exit(1)
	}
}

func handle() {
	if err := utils.Write(Database, "current_ip", []byte(IP().CurIP.String())); err != nil {
		fmt.Println(err)
	}
	// Golang does not support the atexit function from C, so I will write to ip.txt
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

	for {
		// If the CPUSaver is enabled, we will only start a new instance if there is one available
		if Config.CPUSaver {
			if tasks < Config.InstanceCount {
				tasks++
				go HandleServer(IP().GetNext(tasks))
				/*switch Arguments.Mode {
				case "random":
					tasks++
					go HandleServer(utils.RandIP())
				case "order":
					tasks++
					go HandleServer(IP().GetNext(tasks))
				}*/
			}
		} else {
			tasks++
			switch Arguments.Mode {
			case "random":
				go HandleServer(utils.RandIP())
			case "order":
				go HandleServer(IP().GetNext(tasks))
			}
		}
	}
}

func HandleServer(ip net.IP) {
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
