package output

import (
	"github.com/snormore/gowire/message"
	"launchpad.net/tomb"
	"sync"
)

type Outputter interface {
	Push(msg message.Message) error
}

var adapter *Outputter

func Init(l *Outputter) {
	adapter = l
}

func Start(messages chan message.Message, errs chan error, wg *sync.WaitGroup, t *tomb.Tomb) error {
	defer wg.Done()

	for {
		select {
		case <-t.Dying():
			return t.Err()
		case msg := <-messages:
			err := (*adapter).Push(msg)
			if err != nil {
				errs <- err
			}
		}
	}
}
