package config

type Config interface{}

func Register(name string, c Config) error {
	return nil
}

type Configurable struct {
	Config interface{}
}
