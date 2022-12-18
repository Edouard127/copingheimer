package intf

import (
	"bytes"
	"encoding/binary"
)

type Arguments struct {
	Help bool

	// Arguments
	Config    string // -c, --config: Path to the config file
	Mode      string // -m, --mode: Mode to run in (default: "random") (random, order)
	IP        string // -ip : IP address to start from, only used in order mode
	CPUSaver  bool   // -cs, --cpu-saver: Whether to save CPU or not
	Instances int    // -i, --instances: Number of instances to run
	Timeout   int    // -t, --timeout: Timeout for each ping

	// Blacklist
	BlacklistFile string // -bf, --blacklist-file: Path to the blacklist file (default: "blacklist.txt") separated by newlines (e.g. 0.0.0.0-10.0.0.0, 10.0.0.0/8)

	// Database
	Database    string // -d, --database: The type of the database, currently supported: bolt, mongodb
	DatabaseURL string // -du, --database-url: The URL to the database
}

func (a *Arguments) ReadFrom(b []byte) error {
	reader := bytes.NewReader(b)
	if err := binary.Read(reader, binary.LittleEndian, a); err != nil {
		return err
	}
	return nil
}
