package wire

import (
	"launchpad.net/tomb"
	"sync"
)

type Outputter interface {
	Start(t *tomb.Tomb) error
	Push(msg Message) error
}

var (
	adapter  *Outputter
	inputter *Inputter
)

func Init(out *Outputter, in *Inputter) {
	adapter = out
	inputter = in
}

func Start(out *Outputter, numberOfListeners int, messages chan Message, errs chan error, t *tomb.Tomb) {
	go func() {
		err := (*out).Start(t)
		if err != nil {
			errs <- err
		}
	}()

	var outWaits sync.WaitGroup
	outWaits.Add(numberOfListeners)
	for i := 0; i < numberOfListeners; i++ {
		go Listen(messages, errs, &outWaits, t)
	}
}

func Listen(messages chan Message, errs chan error, wg *sync.WaitGroup, t *tomb.Tomb) error {
	defer wg.Done()

	for {
		select {
		case <-t.Dying():
			return t.Err()
		case msg := <-messages:
			err := (*adapter).Push(msg)
			if err != nil {
				errs <- err
			} else {
				(*inputter).FinalizeMessage(msg)
			}
		}
	}
}
