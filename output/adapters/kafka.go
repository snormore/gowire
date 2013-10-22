package kafka_outputter

import (
	"github.com/snormore/gowire"
	"github.com/snormore/gowire/message"
	"sync"
)

type KafkaOutputter struct{}

func (l KafkaOutputter) Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error {
	return nil
}
