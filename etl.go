package etl

import (
	"github.com/snormore/goetl/config"
	"github.com/snormore/goetl/extract"
	"github.com/snormore/goetl/load"
	"github.com/snormore/goetl/message"
	"launchpad.net/tomb"
	"sync"
)

var Config = DefaultConfig

func Init() {
	InitEx(nil)
}

func InitEx(c *EtlConfig) {
	if c != nil {
		Config = *c
	}
	config.Register("etl", Config)
}

func Start(extractor extract.Extractor, loader load.Loader, errs chan error, t *tomb.Tomb) {
	messages := make(chan message.Message, Config.MessagesChannelSize)

	extract.Init(extractor)
	var extractorsWaits sync.WaitGroup
	for i := 0; i < Config.NumberOfExtractors; i++ {
		go extract.Start(messages, errs, &extractorsWaits, t)
	}

	load.Init(loader)
	var loaderWaits sync.WaitGroup
	for i := 0; i < Config.NumberOfExtractors; i++ {
		go load.Start(messages, errs, &loaderWaits, t)
	}

	select {
	case <-t.Dying():
	}
}
