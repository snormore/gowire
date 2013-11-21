package wire

import (
	"github.com/snormore/goconfig"
	"launchpad.net/tomb"
)

var Config = DefaultConfig

func Init() {
	InitEx(nil)
}

func InitEx(c *WireConfig) {
	if c != nil {
		Config = *c
	}
	config.Register("wire", Config)
}

func Start(in *Inputter, out *Outputter, errs chan error, t *tomb.Tomb) {
	messages := make(chan Message, Config.BufferSize)

	Init(in)
	go Start(in, Config.NumberOfInputters, messages, errs, t)

	Init(out, in)
	go Start(out, Config.NumberOfOutputters, messages, errs, t)

	<-t.Dying()
}
