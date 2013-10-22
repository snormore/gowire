package wire

import (
	"github.com/snormore/gowire/config"
	"github.com/snormore/gowire/input"
	"github.com/snormore/gowire/message"
	"github.com/snormore/gowire/output"
	"launchpad.net/tomb"
	"sync"
)

var Config = DefaultConfig

func Init() {
	InitEx(nil)
}

func InitEx(c *WireConfig) {
	if c != nil {
		Config = *c
	}
	config.Register("etl", Config)
}

func Start(in *input.Inputter, out *output.Outputter, errs chan error, t *tomb.Tomb) {
	messages := make(chan message.Message, Config.MessagesChannelSize)

	input.Init(in)
	var inWaits sync.WaitGroup
	inWaits.Add(Config.NumberOfInputters)
	for i := 0; i < Config.NumberOfInputters; i++ {
		go input.Start(messages, errs, &inWaits, t)
	}

	output.Init(out)
	var outWaits sync.WaitGroup
	outWaits.Add(Config.NumberOfOutputters)
	for i := 0; i < Config.NumberOfInputters; i++ {
		go output.Start(messages, errs, &outWaits, t)
	}

	select {
	case <-t.Dying():
	}
}
