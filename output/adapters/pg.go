package pg_outputter

import (
	"github.com/snormore/gowire"
	"github.com/snormore/gowire/message"
	"sync"
)

type PgOutputter struct{}

func (l PgOutputter) Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error {
	return nil
}
