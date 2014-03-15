package wire

type Config struct {
	NumberOfInputters    int `json:"number_of_inputters"`
	NumberOfTransformers int `json:"number_of_transformers"`
	NumberOfOutputters   int `json:"number_of_outputters"`
	BufferSize           int `json:"buffer_size"`
}

var DefaultConfig = Config{
	NumberOfInputters:    10,
	NumberOfTransformers: 10,
	NumberOfOutputters:   10,
	BufferSize:           1024,
}

func NewConfig(rawConfig map[string]interface{}) (*Config, error) {
	config := new(Config)

	if _, ok := rawConfig["number_of_inputters"]; !ok {
		config.NumberOfInputters = DefaultConfig.NumberOfInputters
	} else {
		config.NumberOfInputters = rawConfig["number_of_inputters"].(int)
	}

	if _, ok := rawConfig["number_of_outputters"]; !ok {
		config.NumberOfOutputters = DefaultConfig.NumberOfOutputters
	} else {
		config.NumberOfOutputters = rawConfig["number_of_outputters"].(int)
	}

	if _, ok := rawConfig["buffer_size"]; !ok {
		config.BufferSize = DefaultConfig.BufferSize
	} else {
		config.BufferSize = rawConfig["buffer_size"].(int)
	}

	return config, nil
}
