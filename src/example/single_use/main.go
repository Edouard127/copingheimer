package main

import (
	"edouard127/copingheimer/src/server"
	"edouard127/copingheimer/src/worker"
)

func main() {
	go server.StartServer()
	go worker.StartWorker()
	select {}
}
