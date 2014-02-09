package wire

import (
	"launchpad.net/tomb"
	"sync"
)

type Outputter interface {
	Start(t *tomb.Tomb) error
	Push(msg interface{}) error
	Close() error
}

type output struct {
	out Outputter
	in  Inputter
	t   *tomb.Tomb
}

func newOutput(out Outputter, in Inputter) *output {
	o := output{out, in, new(tomb.Tomb)}
	return &o
}

func (o *output) start(numberOfListeners int, messages chan interface{}, errs chan error) error {

	err := o.out.Start(o.t)
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
		go o.listen(messages, errs, &outWaits)
	}

	return nil
}

func (o *output) listen(messages chan interface{}, errs chan error, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		case <-o.t.Dying():
			return o.t.Err()
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

func (o *output) close() error {
	select {
	case <-o.t.Dying():
	default:
		o.t.Done()
	}
	return o.out.Close()
}
