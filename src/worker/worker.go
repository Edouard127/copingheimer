package worker

import (
	"edouard127/copingheimer/src/intf"
	"edouard127/copingheimer/src/worker/structure"
	"flag"
	"net"
	"os"
)

var (
	Arguments = intf.Arguments{}
	worker    = &structure.Worker{}
)

func StartWorker() {
	flag.BoolVar(&Arguments.Help, "h", false, "Show this help message")
	flag.BoolVar(&Arguments.Help, "help", false, "Show this help message")
	flag.StringVar(&Arguments.Node, "n", "localhost:29229", "The node to connect to")
	flag.StringVar(&Arguments.Node, "node", "localhost:29229", "The node to connect to")
	flag.StringVar(&Arguments.IP, "ip", "109.237.24.0", "IP address to start from")
	flag.IntVar(&Arguments.Instances, "i", 256, "Number of instances to run (default: 1)")
	flag.IntVar(&Arguments.Instances, "instances", 256, "Number of instances to run (default: 1)")
	flag.IntVar(&Arguments.Timeout, "t", 1000, "Timeout for each ping (default: 1000)")
	flag.IntVar(&Arguments.Timeout, "timeout", 1000, "Timeout for each ping (default: 1000)")
	flag.StringVar(&Arguments.BlacklistFile, "bf", "blacklist.txt", "Path to the blacklist file")
	flag.StringVar(&Arguments.BlacklistFile, "blacklist-file", "blacklist.txt", "Path to the blacklist file")
	flag.Parse()

	if Arguments.Help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	worker = structure.NewWorker(&Arguments)
	worker.Connect(net.TCPAddr{IP: net.ParseIP(Arguments.Node), Port: 29229})
	worker.Handle() // This is where the magic happens
}
