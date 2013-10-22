package wire

import (
	"fmt"
	"github.com/snormore/gologger"
	"github.com/snormore/gowire/input"
	"github.com/snormore/gowire/message"
	"github.com/snormore/gowire/output"
	"github.com/stretchr/testify/assert"
	"launchpad.net/tomb"
	"testing"
	"time"
)

type FakeInputter struct {
	count int
}

func (e FakeInputter) Transform(rawMessage interface{}) (message.Message, error) {
	e.count++
	msg := message.Message{string(e.count), rawMessage}
	return msg, nil
}

func (e FakeInputter) Listen() chan interface{} {
	messages := make(chan interface{}, 1024)
	go func() {
		for i := 0; i < 1000; i++ {
			messages <- fmt.Sprintf("Message #%d.", i+1)
			time.Sleep(500 * time.Nanosecond)
		}
	}()
	return messages
}

type FakeOutputter struct{}

func (e FakeOutputter) Push(msg message.Message) error {
	logger.Info("Loader pushing message: %+v", msg)
	return nil
}

func TestStart(t *testing.T) {
	in := FakeInputter{}
	out := FakeOutputter{}
	errs := make(chan error, 1024)
	go func() {
		for err := range errs {
			assert.NoError(t, err)
		}
	}()
	var startTomb tomb.Tomb
	go Start(input.Inputter(in), output.Outputter(out), errs, &startTomb)
	startTomb.Done()
}
