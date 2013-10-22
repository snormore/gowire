package tail_extractor

import (
	"github.com/snormore/goetl"
	"github.com/snormore/goetl/message"
	"sync"
)

type TailExtractor struct{}

func (e TailExtractor) SetTransformer(t Transformer) error {
	return nil
}

func (e TailExtractor) Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error {
	return nil
}
