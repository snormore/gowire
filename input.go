package wire

import (
	"launchpad.net/tomb"
	"sync"
)

type Inputter interface {
	Start(t *tomb.Tomb) error
	Listen() chan interface{}
	Transform(rawMessage interface{}) (interface{}, error)
	FinalizeMessage(msg interface{}) error
	Close() error
}

type input struct {
	in          Inputter
	transformer Transformer
	t           *tomb.Tomb
}

var adapter Inputter

func newInput(in Inputter, transformer Transformer) *input {
	i := input{in, transformer, new(tomb.Tomb)}
	return &i
}

func (i *input) start(numberOfListeners int, messages chan interface{}, errs chan error) error {

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

func (i *input) listen(messages chan interface{}, errs chan error, wg *sync.WaitGroup) error {
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
			msg, err := i.transformer.Transform(rawMsg)
			if err != nil {
				errs <- err
			} else if msg != nil {
				messages <- msg
			}
		}
	}
}

func (i *input) close() error {
	i.t.Done()
	return i.in.Close()
}
