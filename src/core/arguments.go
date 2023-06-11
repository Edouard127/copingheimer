package core

type ClientOption struct {
	Node      string
	Instances int
	Timeout   int
}

type ServerOption struct {
	MongoDB string
	Host    string
	StartIP uint
}
