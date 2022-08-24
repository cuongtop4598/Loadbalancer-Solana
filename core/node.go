package core

type Node struct {
	Endpoint    string
	BlockNumber uint64
	Available   bool
	RPCCounter  int
}
