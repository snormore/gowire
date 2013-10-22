package load

import (
	"github.com/snormore/goetl"
	"github.com/snormore/goetl/message"
	"sync"
)

type Loader interface {
	Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error
}
