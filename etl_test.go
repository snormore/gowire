package etl

import (
	"fmt"
	"github.com/snormore/goetl/extract"
	"github.com/snormore/goetl/load"
	"github.com/snormore/goetl/message"
	"github.com/snormore/gologger"
	"github.com/stretchr/testify/assert"
	"launchpad.net/tomb"
	"testing"
	"time"
)

type FakeExtractor struct {
	count int
}

func (e FakeExtractor) Transform(rawMessage interface{}) (message.Message, error) {
	e.count++
	msg := message.Message{string(e.count), rawMessage}
	return msg, nil
}

func (e FakeExtractor) Listen() chan interface{} {
	messages := make(chan interface{}, 1024)
	go func() {
		for i := 0; i < 1000; i++ {
			messages <- fmt.Sprintf("Message #%d.", i+1)
			time.Sleep(500 * time.Nanosecond)
		}
	}()
	return messages
}

type FakeLoader struct{}

func (e FakeLoader) Push(msg message.Message) error {
	logger.Info("Loader pushing message: %+v", msg)
	return nil
}

func TestStart(t *testing.T) {
	extractor := FakeExtractor{}
	loader := FakeLoader{}
	errs := make(chan error, 1024)
	go func() {
		for err := range errs {
			assert.NoError(t, err)
		}
	}()
	var startTomb tomb.Tomb
	Start(extract.Extractor(extractor), load.Loader(loader), errs, &startTomb)
	startTomb.Done()
}
