package config

type AbstractConfig map[string]interface{}

var Config AbstractConfig{}

func Register(name string, c AbstractConfig) error {
	return nil
}
