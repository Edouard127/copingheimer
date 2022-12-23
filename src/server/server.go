package server

import (
	"edouard127/copingheimer/src/net/packet"
	"fmt"
	"net"
)

var listn = &packet.Listener{}

func StartServer(mongo string) {
	server := packet.NewServer(mongo)
	var err error
	taddr := net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 29229,
	}
	if listn, err = packet.ListenProvider(taddr); err != nil {
		panic(err)
	}

	fmt.Println("Server started")

	server.Events.AddListener(
		packet.PacketHandler{ID: packet.SPacketLogin, F: server.EventHandlers.ConfirmRegistration},
		packet.PacketHandler{ID: packet.SPacketIP, F: server.EventHandlers.ConfirmIP},
	)

	go server.Dashboard.Start()

	for {
		if conn, err := listn.Accept(); err != nil {
			server.Remove(conn)
		} else {
			server.Add(conn)
			fmt.Println("New connection")
			go func() {
				var p packet.Packet
				for {
					if err := conn.ReadPacket(&p); err == nil {
						if err := server.HandlePacket(conn, p); err != nil {
							fmt.Println(err)
						}
					}
				}
			}()
		}
	}
}
