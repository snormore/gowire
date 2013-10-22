package extract

import (
	"github.com/snormore/goetl/message"
	"launchpad.net/tomb"
	"sync"
)

type Extractor interface {
	Listen() chan interface{}
	Transform(message interface{}) (message.Message, error)
}

var adapter Extractor

func Init(e Extractor) {
	adapter = e
}

func Start(messages chan message.Message, errs chan error, wg *sync.WaitGroup, t *tomb.Tomb) error {
	defer wg.Done()

	for {
		select {
		case <-t.Dying():
			return t.Err()
		case rawMsg := <-adapter.Listen():
			msg, err := adapter.Transform(rawMsg)
			if err == nil {
				messages <- msg
			} else {
				errs <- err
			}
		}
	}
}
