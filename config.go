package etl

type EtlConfig struct {
	NumberOfExtractors  int
	NumberOfLoaders     int
	MessagesChannelSize int
}

var DefaultConfig = EtlConfig{
	NumberOfExtractors:  10,
	NumberOfLoaders:     10,
	MessagesChannelSize: 1024,
}
