package wire

import (
	"github.com/stretchr/testify/assert"
	"launchpad.net/tomb"
	"strings"
	"testing"
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

func sampleEventLog() string {
	return strings.Trim(SampleEventLog, " \n\r\t")
}

func sampleEventLogEntries() []string {
	sampleEntries := strings.Split(sampleEventLog(), "\n")
	entries := make([]string, 0, len(sampleEntries))
	for _, entry := range sampleEntries {
		if strings.Trim(entry, " \n\r\t") != "" {
			entries = append(entries, entry)
		}
	}
	return entries
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
	in := NewFakeInputter()
	var inTomb tomb.Tomb
	go in.Start(&inTomb)

	sampleLogEntries := sampleEventLogEntries()
	go in.PushAll(sampleLogEntries)

	outMessages := make(chan Message, 1024)
	out := FakeOutputter{outMessages}

	errs := consumeAndCheckErrors(t)

	var startTomb tomb.Tomb
	inputter := Inputter(in)
	outputter := Outputter(out)
	go Start(&inputter, &outputter, errs, &startTomb)

	i := 0
	for _ = range out.Messages {
		i++
		if i == FakeInputterMessageCount {
			startTomb.Killf("TestStartWithMocks: ending...")
			break
		}
	}
	assert.Equal(t, FakeInputterMessageCount, i)
	startTomb.Wait()
}
