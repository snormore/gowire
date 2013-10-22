package extract

import (
	"github.com/snormore/goetl"
	"github.com/snormore/goetl/message"
	"sync"
)

type Extractor interface {
	SetTransformer(t Transformer) error
	Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error
}
