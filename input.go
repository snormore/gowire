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
	t  *tomb.Tomb
}

var adapter Inputter

func newInput(in Inputter) *input {
	i := input{in, new(tomb.Tomb)}
	return &i
}

func (i *input) start(numberOfListeners int, messages chan Message, errs chan error) error {

	err := i.in.Start(i.t)
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
		go i.listen(messages, errs, &inWaits)
	}

	return nil
}

func (i *input) listen(messages chan Message, errs chan error, wg *sync.WaitGroup) error {
	defer func() {
		wg.Done()
		select {
		case <-i.t.Dead():
		default:
			i.close()
		}
	}()

	for {
		select {
		case <-i.t.Dying():
			return i.t.Err()
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

func (i *input) close() error {
	i.t.Done()
	return i.in.Close()
}
