package wire

import (
	"launchpad.net/tomb"
)

type Wire struct {
	in  Inputter
	out Outputter

	Config *WireConfig
}

func New(config *WireConfig) *Wire {
	w := new(Wire)
	if config == nil {
		w.Config = &DefaultConfig
	} else {
		w.Config = config
	}
	return w
}

func (w *Wire) Start(in Inputter, out Outputter, errs chan error, t *tomb.Tomb) error {
	messages := make(chan Message, w.Config.BufferSize)

	i := newInput(in)
	if err := i.start(w.Config.NumberOfInputters, messages, errs, t); err != nil {
		return err
	}

	o := newOutput(out, in)
	if err := o.start(w.Config.NumberOfOutputters, messages, errs, t); err != nil {
		return err
	}

	return nil
}
