package sky_loader

import (
	"github.com/snormore/goetl"
	"github.com/snormore/goetl/message"
	"sync"
)

type SkyLoader struct{}

func (l SkyLoader) Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error {
	return nil
}
