package wire

import (
	"launchpad.net/tomb"
	"sync"
)

type Outputter interface {
	Start(t *tomb.Tomb) error
	Push(msg Message) error
	Close() error
}

type output struct {
	out Outputter
	in  Inputter
}

func newOutput(out Outputter, in Inputter) *output {
	o := output{out, in}
	return &o
}

func (o *output) start(numberOfListeners int, messages chan Message, errs chan error, t *tomb.Tomb) error {

	err := o.out.Start(t)
	if err != nil {
		return err
	}

	go func() {
		if err != nil {
			errs <- err
		}
	}()

	var outWaits sync.WaitGroup
	outWaits.Add(numberOfListeners)
	for i := 0; i < numberOfListeners; i++ {
		go o.listen(messages, errs, &outWaits, t)
	}

	return nil
}

func (o *output) listen(messages chan Message, errs chan error, wg *sync.WaitGroup, t *tomb.Tomb) error {
	defer wg.Done()

	for {
		select {
		case <-t.Dying():
			return t.Err()
		case msg := <-messages:
			err := o.out.Push(msg)
			if err != nil {
				errs <- err
			} else {
				o.in.FinalizeMessage(msg)
			}
		}
	}
}
