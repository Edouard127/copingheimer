package main

import (
	"edouard127/copingheimer/src/server"
	"edouard127/copingheimer/src/worker"
	"flag"
)

func main() {
	// Parse flags
	types := flag.String("type", "server", "Type of the program to run (server or worker)")
	flag.Parse()
	if *types == "server" {
		server.StartServer()
	} else if *types == "worker" {
		worker.StartWorker()
	}
}
