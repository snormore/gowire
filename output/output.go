package output

import (
	"github.com/snormore/gowire/input"
	"github.com/snormore/gowire/message"
	"launchpad.net/tomb"
	"sync"
)

type Outputter interface {
	Start() error
	Push(msg Message) error
	Close() error
}

var (
	adapter  Outputter
	inputter input.Inputter
)

func Init(out Outputter, in input.Inputter) {
	adapter = out
	inputter = in
}

func Start(out Outputter, numberOfListeners int, messages chan Message, t *tomb.Tomb) error {
	go func() {
		err := out.Start(t)
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
			err := adapter.Push(msg)
			if err != nil {
				errs <- err
			} else {
				inputter.FinalizeMessage(msg)
			}
		}
	}
}
