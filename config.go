package wire

type WireConfig struct {
	NumberOfInputters   int
	NumberOfOutputters  int
	MessagesChannelSize int
}

var DefaultConfig = WireConfig{
	NumberOfInputters:   10,
	NumberOfOutputters:  10,
	MessagesChannelSize: 1024,
}
