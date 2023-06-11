package main

import (
	"edouard127/copingheimer/src/core"
	"edouard127/copingheimer/src/net"
	"edouard127/copingheimer/src/net/packet"
	"edouard127/copingheimer/src/worker"
	"flag"
)

func main() {
	options := &core.ClientOption{
		Node:      "127.0.0.1:29969",
		Instances: 256,
		Timeout:   4000,
	}

	flag.StringVar(&options.Node, "node", options.Node, "Server node")
	flag.IntVar(&options.Instances, "instances", options.Instances, "Number of instances")
	flag.IntVar(&options.Timeout, "timeout", options.Timeout, "Timeout in milliseconds")

	w := worker.NewWorker(options)

	w.Events.AddListener(
		packet.Handler[worker.Worker]{ID: net.CPacketKeepAlive, F: worker.KeepAlive},
		packet.Handler[worker.Worker]{ID: net.CPacketOffsetIP, F: worker.OffsetIP},
		packet.Handler[worker.Worker]{ID: net.CPacketSignal, F: worker.Signal},
	)

	w.Handle() // This is a blocking call
}
