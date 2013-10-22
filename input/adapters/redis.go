package input_adapter

import (
	"github.com/snormore/gowire/message"
	"launchpad.net/tomb"
)

type RedisInputter struct{}

type RedisConfig struct {
	Host string
	Port uint32
}

func NewRedisInputter() *RedisInputter {
	in := RedisInputter{}
	return &in
}

func (in RedisInputter) Start(config interface{}, t *tomb.Tomb) error {
	return nil
}

func (in RedisInputter) Transform(rawMessage interface{}) (message.Message, error) {
	return message.Message{}, nil
}

func (in RedisInputter) Listen() chan interface{} {
	return nil
}

func (in RedisInputter) FinalizeMessage(msg message.Message) error {
	return nil
}
