package net

const (
	SPacketLogin = iota
	SPacketKeepAlive
	SPacketServer
	SPacketIP
)

const (
	CPacketKeepAlive = iota
	CPacketOffsetIP
	CPacketSignal
)
