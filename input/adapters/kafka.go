package input_adapter

import (
	"github.com/snormore/gowire/message"
)

type KafkaConfig struct{}

type KafkaInputter struct{}

func (in KafkaInputter) Transform(rawMessage interface{}) (message.Message, error) {
	return message.Message{"", rawMessage}, nil
}

func (in KafkaInputter) Listen() chan interface{} {
	return nil
}
