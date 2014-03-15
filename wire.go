package wire

type Wire struct {
	in  *input
	out *output
	tr  *transform

	Config *Config
}

func New(config *Config) *Wire {
	w := new(Wire)
	if config == nil {
		w.Config = &DefaultConfig
	} else {
		w.Config = config
	}
	if w.Config.NumberOfInputters == 0 {
		w.Config.NumberOfInputters = 1
	}
	if w.Config.NumberOfTransformers == 0 {
		w.Config.NumberOfTransformers = 1
	}
	if w.Config.NumberOfOutputters == 0 {
		w.Config.NumberOfOutputters = 1
	}
	return w
}

func (w *Wire) Start(in Inputter, out Outputter, transformer Transformer, errs chan error) error {
	in_messages := make(chan interface{}, w.Config.BufferSize)
	out_messages := make(chan interface{}, w.Config.BufferSize)

	w.in = newInput(in)
	if err := w.in.start(w.Config.NumberOfInputters, in_messages, errs); err != nil {
		return err
	}

	w.tr = newTransform(transformer)
	if err := w.tr.start(w.Config.NumberOfTransformers, in_messages, out_messages, errs); err != nil {
		return err
	}

	w.out = newOutput(out, in)
	if err := w.out.start(w.Config.NumberOfOutputters, out_messages, errs); err != nil {
		return err
	}

	return nil
}

func (w *Wire) Close() error {
	if err := w.in.close(); err != nil {
		return err
	}
	return w.out.close()
}
