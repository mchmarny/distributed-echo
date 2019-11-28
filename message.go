package main

// BroadcastMessage represents the EchoNodes this node should echo
type BroadcastMessage struct {
	Targets []*EchoNode `yaml:"targets"`
}

// EchoNode represents single broadcast target
type EchoNode struct {
	Region string `yaml:"region"`
	URL    string `yaml:"url"`
}

// EchoMessage is the message sent from one node to another
type EchoMessage struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
	Sent int64  `yaml:"send"`
}
