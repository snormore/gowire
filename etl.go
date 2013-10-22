package etl

import (
	"github.com/snormore/goetl/extract"
	"github.com/snormore/goetl/load"
	"github.com/snormore/goetl/logger"
	"github.com/snormore/goetl/message"
	"github.com/snormore/goetl/transform"
	"sync"
)

var Config = DefaultConfig

func Init() {
	InitEx(nil)
}

func InitEx(c EtlConfig) {
	if c != nil {
		Config = c
	}
	config.Register(Config)
}

func Start(extractor extract.Extractor, transformer transform.Transformer, loader load.Loader, t *tomb.Tomb) {
	messages := make(chan message.Message)

	var extractorsWaits sync.WaitGroup
	for i := 0; i < Config.NumberOfExtractors; i++ {
		go extractor.Start(messages, &extractorsWaits, t)
	}

	var loaderWaits sync.WaitGroup
	for i := 0; i < Config.NumberOfExtractors; i++ {
		go loader.Start(messages, &loaderWaits, t)
	}

	select {
	case <-t.Dying():
	}

}
