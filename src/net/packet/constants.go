package packet

const (
	CPacketLogin = iota + 1 // Returns the IP start from the worker's ID
	CPacketUpdate
	CPacketIP // Puts the worker on wait until all workers have sent their IP
)

const (
	SPacketLogin  = iota + 1 // Sends the worker's ID
	SPacketIP                // Sends the worker's current IP
	SPacketResult            // Sends the worker's result
)
