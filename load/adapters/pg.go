package pg_loader

import (
	"github.com/snormore/goetl"
	"github.com/snormore/goetl/message"
	"sync"
)

type PgLoader struct{}

func (l PgLoader) Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error {
	return nil
}
