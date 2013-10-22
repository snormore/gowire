package config

type AbstractConfig map[string]interface{}

func Register(name string, c interface{}) error {
	return nil
}
