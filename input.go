package wire

import (
	"launchpad.net/tomb"
	"sync"
)

type Inputter interface {
	Start(t *tomb.Tomb) error
	Listen() chan interface{}
	Transform(rawMessage interface{}) (Message, error)
	FinalizeMessage(msg Message) error
	Close() error
}

type input struct {
	in Inputter
}

var adapter Inputter

func newInput(in Inputter) *input {
	i := input{in}
	return &i
}

func (i *input) start(numberOfListeners int, messages chan Message, errs chan error, t *tomb.Tomb) error {

	err := i.in.Start(t)
	if err != nil {
		return err
	}

	go func() {
		if err != nil {
			errs <- err
		}
	}()

	var inWaits sync.WaitGroup
	inWaits.Add(numberOfListeners)
	for j := 0; j < numberOfListeners; j++ {
		go i.listen(messages, errs, &inWaits, t)
	}

	return nil
}

func (i *input) listen(messages chan Message, errs chan error, wg *sync.WaitGroup, t *tomb.Tomb) error {
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
		case rawMsg := <-i.in.Listen():
			msg, err := i.in.Transform(rawMsg)
			if err == nil {
				messages <- msg
			} else {
				errs <- err
			}
		}
	}
}
