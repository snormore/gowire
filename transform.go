package wire

import (
	"launchpad.net/tomb"
	"sync"
)

type Transformer interface {
	Start(t *tomb.Tomb) error
	Transform(msg interface{}) (interface{}, error)
	Close() error
}

type transform struct {
	tr Transformer
	t  *tomb.Tomb
}

func newTransform(tr Transformer) *transform {
	t := transform{tr, new(tomb.Tomb)}
	return &t
}

func (tr *transform) start(numberOfListeners int, in_messages chan interface{}, out_messages chan interface{}, errs chan error) error {

	err := tr.tr.Start(tr.t)
	if err != nil {
		return err
	}

	go func() {
		if err != nil {
			errs <- err
		}
	}()

	var trWaits sync.WaitGroup
	trWaits.Add(numberOfListeners)
	for j := 0; j < numberOfListeners; j++ {
		go tr.listen(in_messages, out_messages, errs, &trWaits)
	}

	return nil
}

func (tr *transform) listen(in_messages chan interface{}, out_messages chan interface{}, errs chan error, wg *sync.WaitGroup) error {
	defer func() {
		wg.Done()
		select {
		case <-tr.t.Dead():
		default:
		}
	}()

	for {
		select {
		case <-tr.t.Dying():
			return tr.t.Err()
		case rawMsg := <-in_messages:
			msg, err := tr.tr.Transform(rawMsg)
			if err != nil {
				errs <- err
			} else if msg != nil {
				out_messages <- msg
			}
		}
	}
}
