package worker

import (
	"edouard127/copingheimer/src/intf"
	"edouard127/copingheimer/src/worker/structure"
	"net"
)

var (
	worker = &structure.Worker{}
)

func StartWorker(args *intf.Arguments) {
	worker = structure.NewWorker(args)
	worker.Connect(net.TCPAddr{IP: net.ParseIP(args.Node), Port: 29229})
	worker.Handle() // This is where the magic happens
}
