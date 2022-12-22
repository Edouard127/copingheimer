package intf

import (
	"bytes"
	"encoding/binary"
)

type Arguments struct {
	Help bool

	// Arguments
	Node      string // -n, --node: Main node to connect to (default: "localhost:29229")
	IP        string // -ip : IP address to start from, only used in order mode
	Instances int    // -i, --instances: Number of instances to run
	Timeout   int    // -t, --timeout: Timeout for each ping
	Hosting   bool   // -h, --hosting: Whether to scan hosters IPs or not

	// Blacklist
	BlacklistFile string // -bf, --blacklist-file: Path to the blacklist file (default: "blacklist.txt") separated by newlines (e.g. 0.0.0.0-10.0.0.0, 10.0.0.0/8)
}

func (a *Arguments) ReadFrom(b []byte) error {
	reader := bytes.NewReader(b)
	if err := binary.Read(reader, binary.LittleEndian, a); err != nil {
		return err
	}
	return nil
}
