package wire

import (
	"fmt"
	"github.com/snormore/gologger"
	"github.com/snormore/gowire/input"
	"github.com/snormore/gowire/message"
	"github.com/snormore/gowire/output"
	"github.com/stretchr/testify/assert"
	"launchpad.net/tomb"
	"strings"
	"testing"
	"time"
)

const FakeInputterMessageCount = 10

func init() {
	logger.Verbosity = 0
}

type FakeInputter struct {
	Messages chan interface{}
	Count    int
}

func (in FakeInputter) Transform(rawMessage interface{}) (message.Message, error) {
	msg := message.Message{strings.Split(rawMessage.(string), "#")[1], rawMessage}
	return msg, nil
}

func (in FakeInputter) Listen() chan interface{} {
	return in.Messages
}

func (in *FakeInputter) pushRawMessages() {
	for in.Count = 0; in.Count < FakeInputterMessageCount; in.Count++ {
		rawMsg := fmt.Sprintf("Message #%d", in.Count+1)
		logger.Debug("Input: %s", rawMsg)
		in.Messages <- rawMsg
		time.Sleep(500 * time.Nanosecond)
	}
}

type FakeOutputter struct {
	Messages chan message.Message
}

func (out FakeOutputter) Push(msg message.Message) error {
	out.Messages <- msg
	return nil
}

func consumeAndCheckErrors(t *testing.T) chan error {
	errs := make(chan error, 1024)
	go func() {
		for err := range errs {
			assert.NoError(t, err)
		}
	}()
	return errs
}

func TestStart(t *testing.T) {
	in := FakeInputter{}
	in.Messages = make(chan interface{}, 1024)
	go in.pushRawMessages()

	outputs := make(chan message.Message, 1024)
	out := FakeOutputter{outputs}

	errs := consumeAndCheckErrors(t)

	var startTomb tomb.Tomb
	inputter := input.Inputter(in)
	outputter := output.Outputter(out)
	go Start(&inputter, &outputter, errs, &startTomb)

	i := 0
	for msg := range out.Messages {
		i++
		logger.Debug("Output: %+v", msg)
		assert.Equal(t, fmt.Sprintf("%d", i), msg.Id)
		if i == FakeInputterMessageCount {
			startTomb.Killf("TestStart: ending...")
			break
		}
	}
	assert.Equal(t, FakeInputterMessageCount, in.Count)
	assert.Equal(t, FakeInputterMessageCount, i)
	startTomb.Wait()
}
