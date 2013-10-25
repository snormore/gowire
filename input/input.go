package input

import (
	"github.com/snormore/gowire/message"
	"launchpad.net/tomb"
	"sync"
)

type Inputter interface {
	Listen() chan interface{}
	Transform(rawMessage interface{}) (message.Message, error)
	FinalizeMessage(msg message.Message) error
	Start(t *tomb.Tomb) error
}

var adapter *Inputter

func Init(e *Inputter) {
	adapter = e
}

func Start(in *Inputter, numberOfListeners int, messages chan message.Message, errs chan error, t *tomb.Tomb) {
	go func() {
		err := (*in).Start(t)
		if err != nil {
			errs <- err
		}
	}()

	var inWaits sync.WaitGroup
	inWaits.Add(numberOfListeners)
	for i := 0; i < numberOfListeners; i++ {
		go Listen(messages, errs, &inWaits, t)
	}
}

func Listen(messages chan message.Message, errs chan error, wg *sync.WaitGroup, t *tomb.Tomb) error {
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
