package wire

import (
	"github.com/snormore/goconfig"
	"github.com/snormore/gowire/input"
	"github.com/snormore/gowire/message"
	"github.com/snormore/gowire/output"
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

func Start(in *input.Inputter, out *output.Outputter, errs chan error, t *tomb.Tomb) {
	messages := make(chan message.Message, Config.MessagesChannelSize)

	input.Init(in)
	go input.Start(in, Config.NumberOfInputters, messages, errs, t)

	output.Init(out)
	go output.Start(out, Config.NumberOfOutputters, messages, errs, t)

	<-t.Dying()
}
