package wire

import (
	"fmt"
	"github.com/snormore/gologger"
	"github.com/snormore/gotail"
	"github.com/snormore/gowire-adapters/tail"
	"github.com/snormore/gowire/input"
	"github.com/snormore/gowire/message"
	"github.com/snormore/gowire/output"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"launchpad.net/tomb"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	FakeInputterMessageCount = 10
	TestTempDir              = "./tmp"
	SampleEventLog           = `
			{"event_id":"110"}
			{"event_id":"111"}
			{"event_id":"112"}
			{"event_id":"113"}
			{"event_id":"114"}
			{"event_id":"115"}
			{"event_id":"116"}
			{"event_id":"117"}
			{"event_id":"118"}
			{"event_id":"119"}
	`
)

func init() {
	logger.Verbosity = 2
	tail.ScriptPath = "../gotail/sbin"
}

type FakeInputter struct {
	Messages chan interface{}
}

func (in FakeInputter) Transform(rawMessage interface{}) (message.Message, error) {
	msg := message.Message{strings.Split(rawMessage.(string), "#")[1], rawMessage}
	return msg, nil
}

func (in FakeInputter) Listen() chan interface{} {
	return in.Messages
}

func (in FakeInputter) Start(t *tomb.Tomb) error {
	for i := 0; i < FakeInputterMessageCount; i++ {
		rawMsg := fmt.Sprintf("Message #%d", i+1)
		logger.Debug("Input: %s", rawMsg)
		in.Messages <- rawMsg
		time.Sleep(500 * time.Nanosecond)
	}
	return nil
}

func (in FakeInputter) FinalizeMessage(msg message.Message) error {
	return nil
}

type FakeOutputter struct {
	Messages chan message.Message
}

func (out FakeOutputter) Push(msg message.Message) error {
	out.Messages <- msg
	return nil
}

func sampleEventLog() string {
	return strings.Trim(SampleEventLog, " \n\r\t")
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

func TestStartWithMocks(t *testing.T) {
	in := FakeInputter{}
	in.Messages = make(chan interface{}, 1024)
	var inTomb tomb.Tomb
	go in.Start(&inTomb)

	outMessages := make(chan message.Message, 1024)
	out := FakeOutputter{outMessages}

	errs := consumeAndCheckErrors(t)

	var startTomb tomb.Tomb
	inputter := input.Inputter(in)
	outputter := output.Outputter(out)
	go Start(&inputter, &outputter, errs, &startTomb)

	i := 0
	for msg := range out.Messages {
		i++
		logger.Debug("Output: %+v", msg)
		if i == FakeInputterMessageCount {
			startTomb.Killf("TestStartWithMocks: ending...")
			break
		}
	}
	assert.Equal(t, FakeInputterMessageCount, i)
	startTomb.Wait()
}

func createAndPushToTempFile(log string) (*os.File, error) {
	file, err := ioutil.TempFile("./tmp/", "tailer-test")
	if err != nil {
		return file, err
	}
	lines := strings.Split(strings.Trim(log, " \n\r\t"), "\n")
	go func() {
		for _, line := range lines {
			file.WriteString(fmt.Sprintf("%s\n", strings.Trim(line, " \t\n\r")))
			time.Sleep(100 * time.Nanosecond)
		}
	}()
	return file, err
}

func TestStartWithTailerAndMockOutput(t *testing.T) {
	inFile, err := createAndPushToTempFile(sampleEventLog())
	assert.NoError(t, err)
	inConfig := tail_adapter.TailConfig{
		FilePath:   inFile.Name(),
		StartEvent: "111",
	}
	in := tail_adapter.NewTailInputter(inConfig)

	outMessages := make(chan message.Message, 1024)
	out := FakeOutputter{outMessages}

	errs := consumeAndCheckErrors(t)

	var startTomb tomb.Tomb
	inputter := input.Inputter(in)
	outputter := output.Outputter(out)
	go Start(&inputter, &outputter, errs, &startTomb)

	i := 0
	for msg := range out.Messages {
		logger.Debug("Output: %+v", msg)
		assert.Equal(t, fmt.Sprintf("11%d", i+2), msg.Id)
		i++
		if i == len(strings.Split(sampleEventLog(), "\n"))-2 {
			startTomb.Killf("TestStartWithTailerAndMockOutput: ending...")
			break
		}
	}
	assert.Equal(t, FakeInputterMessageCount-2, i)
	startTomb.Wait()
}

func TestStartWithTailerStartLastEventAndMockOutput(t *testing.T) {
	inFile, err := createAndPushToTempFile(sampleEventLog())
	assert.NoError(t, err)
	inConfig := tail_adapter.TailConfig{
		FilePath:   inFile.Name(),
		StartEvent: "119",
	}
	in := tail_adapter.NewTailInputter(inConfig)

	outMessages := make(chan message.Message, 1024)
	out := FakeOutputter{outMessages}

	errs := consumeAndCheckErrors(t)

	var startTomb tomb.Tomb
	inputter := input.Inputter(in)
	outputter := output.Outputter(out)
	go Start(&inputter, &outputter, errs, &startTomb)

	time.Sleep(50 * time.Millisecond)
	anyMessage := false
	select {
	case <-out.Messages:
		anyMessage = true
	default:
	}
	assert.False(t, anyMessage, "Messages were found on channel where it should be empty.")
	startTomb.Killf("TestStartWithTailerStartLastEventAndMockOutput: ending...")
	startTomb.Wait()
}

func TestStartWithMockInputAndSkyOutput(t *testing.T) {
	inFile, err := createAndPushToTempFile(sampleEventLog())
	assert.NoError(t, err)
	inConfig := tail_adapter.TailConfig{
		FilePath:   inFile.Name(),
		StartEvent: "111",
	}
	in := tail_adapter.NewTailInputter(inConfig)

	outMessages := make(chan message.Message, 1024)
	out := FakeOutputter{outMessages}

	errs := consumeAndCheckErrors(t)

	var startTomb tomb.Tomb
	inputter := input.Inputter(in)
	outputter := output.Outputter(out)
	go Start(&inputter, &outputter, errs, &startTomb)

	i := 0
	for msg := range out.Messages {
		logger.Debug("Output: %+v", msg)
		assert.Equal(t, fmt.Sprintf("11%d", i+2), msg.Id)
		i++
		if i == len(strings.Split(sampleEventLog(), "\n"))-2 {
			startTomb.Killf("TestStartWithTailerAndMockOutput: ending...")
			break
		}
	}
	assert.Equal(t, FakeInputterMessageCount-2, i)
	startTomb.Wait()
}
