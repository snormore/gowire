package input

import (
	"github.com/snormore/gowire/message"
	"launchpad.net/tomb"
	"sync"
)

type Inputter interface {
	Listen() chan interface{}
	Transform(rawMessage interface{}) (message.Message, error)
}

var adapter *Inputter

func Init(e *Inputter) {
	adapter = e
}

func Start(messages chan message.Message, errs chan error, wg *sync.WaitGroup, t *tomb.Tomb) error {
	defer func() {
		wg.Done()
		select {
		case <-t.Dead():
		default:
			t.Done()
		}
	}()

	for {
		select {
		case <-t.Dying():
			return t.Err()
		case rawMsg := <-(*adapter).Listen():
			msg, err := (*adapter).Transform(rawMsg)
			if err == nil {
				messages <- msg
			} else {
				errs <- err
			}
		}
	}
}
