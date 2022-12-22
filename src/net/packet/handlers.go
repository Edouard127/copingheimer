package packet

import (
	"fmt"
	"net"
)

type EventHandlers struct{}

func (e *EventHandlers) ConfirmRegistration(s *Server, c *Conn, p Packet) error {
	var id int32

	if err := p.Scan(&id); err != nil {
		fmt.Println(err)
	}

	ip := c.Socket.RemoteAddr().String()
	parsedIP, _, err := net.SplitHostPort(ip)
	if err != nil {
		return err
	}
	IP := net.ParseIP(parsedIP)

	if err := c.WritePacket(
		Marshal(
			CPacketLogin,
			id, // TODO: Make requests to the other clients
			s.Dashboard.CreateNewUser(IP),
		),
	); err != nil {
		return err
	}
	fmt.Println("Registered", len(s.clients), "workers")
	return nil
}

func (e *EventHandlers) ConfirmIP(s *Server, c *Conn, p Packet) error {
	var ip int32

	if err := p.Scan(&ip); err != nil {
		return err
	}

	var mustRelease bool
	for _, v := range s.clients {
		if v.states.Has(Wait) {
			mustRelease = true
			v.states.Set(Wait, false)
		} else {
			mustRelease = false
		}
	}

	if ip >= c.lastState+c.id && len(s.clients) > 1 {
		c.lastState = ip
		c.states.Set(Wait, true)

		if mustRelease {
			fmt.Println("All clients are on hold, releasing...")
			// Release all clients
			for _, client := range s.clients {
				c.states.Set(Wait, false)
				if err := client.WritePacket(
					Marshal(
						CPacketIP,
						ToInt(c.states.Get(Wait)),
					),
				); err != nil {
					return err
				}
			}
		} // We do that there so if there's an error, it will still release the clients before setting workers on hold

		// Put on hold
		if err := c.WritePacket(
			Marshal(
				CPacketIP,
				ToInt(c.states.Get(Wait)),
			),
		); err != nil {
			return err
		}
	}

	return nil
}
