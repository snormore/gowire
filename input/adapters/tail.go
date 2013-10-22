package tail_inputter

import (
	"github.com/snormore/gowire"
	"github.com/snormore/gowire/message"
	"sync"
)

type TailInputter struct{}

func (e TailInputter) SetTransformer(t Transformer) error {
	return nil
}

func (e TailInputter) Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error {
	return nil
}
