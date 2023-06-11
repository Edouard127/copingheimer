package server

import (
	"context"
	"edouard127/copingheimer/src/core"
	"edouard127/copingheimer/src/net"
	"edouard127/copingheimer/src/net/packet"
	"edouard127/copingheimer/src/server"
	"flag"
)

func main() {
	options := &core.ServerOption{
		MongoDB: "mongodb://localhost:27017",
		Host:    "0.0.0.0",
		StartIP: core.IPInt(core.DefaultIP),
	}

	flag.StringVar(&options.MongoDB, "mongo", options.MongoDB, "MongoDB connection string")
	flag.StringVar(&options.Host, "host", options.Host, "Server host")
	flag.UintVar(&options.StartIP, "ip", options.StartIP, "Starting IP")

	s := server.NewServer(options)

	s.Events.AddListener(
		packet.Handler[server.ServerClient]{ID: net.SPacketLogin, F: s.HandleLogin},
		//packet.Handler[server.ServerClient]{ID: net.SPacketKeepAlive, F: s.HandleLogin},
		packet.Handler[server.ServerClient]{ID: net.SPacketServer, F: s.HandleServer},
		packet.Handler[server.ServerClient]{ID: net.SPacketIP, F: s.HandleIP},
	)

	go s.ClientManager.Run(context.Background())
	s.Handle() // This is a blocking call
}
