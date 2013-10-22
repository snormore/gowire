package etl

type EtlConfig struct {
	NumberOfExtractors int
	NumberofLoaders    int
}

var DefaultConfig = EtlConfig{
	NumberOfExtractors: 10,
	NumberOfLoaders:    10,
	MessageChannelSize: 1024,
}
