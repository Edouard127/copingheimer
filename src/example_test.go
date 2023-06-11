package main

import (
	"context"
	"edouard127/copingheimer/src/core"
	"edouard127/copingheimer/src/net"
	"edouard127/copingheimer/src/net/packet"
	"edouard127/copingheimer/src/server"
	"edouard127/copingheimer/src/worker"
	"testing"
)

func TestClientServerInteraction(t *testing.T) {
	s := server.NewServer(&core.ServerOption{
		MongoDB: "mongodb://localhost:27017",
		Host:    "127.0.0.1",
		StartIP: core.IPInt(core.DefaultIP),
	})

	s.Events.AddListener(
		packet.Handler[server.ServerClient]{ID: net.SPacketLogin, F: s.HandleLogin},
		//packet.Handler[server.ServerClient]{ID: net.SPacketKeepAlive, F: s.HandleLogin},
		packet.Handler[server.ServerClient]{ID: net.SPacketServer, F: s.HandleServer},
		packet.Handler[server.ServerClient]{ID: net.SPacketIP, F: s.HandleIP},
	)

	go s.ClientManager.Run(context.Background())
	go s.Handle()

	w := worker.NewWorker(&core.ClientOption{
		Node:      "127.0.0.1:29969",
		Instances: 256,
		Timeout:   4000,
	})

	w.Events.AddListener(
		packet.Handler[worker.Worker]{ID: net.CPacketKeepAlive, F: worker.KeepAlive},
		packet.Handler[worker.Worker]{ID: net.CPacketOffsetIP, F: worker.OffsetIP},
		packet.Handler[worker.Worker]{ID: net.CPacketSignal, F: worker.Signal},
	)

	go w.Handle()

	for {
	}
}
