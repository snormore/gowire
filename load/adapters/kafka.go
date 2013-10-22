package kafka_loader

import (
	"github.com/snormore/goetl"
	"github.com/snormore/goetl/message"
	"sync"
)

type KafkaLoader struct{}

func (l KafkaLoader) Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error {
	return nil
}
