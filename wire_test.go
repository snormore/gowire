package wire

import (
	"github.com/stretchr/testify/assert"
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
	sampleLogEntries := sampleEventLogEntries()

	outMessages := make(chan interface{}, 1024)
	out := &FakeOutputter{outMessages}

	w := New(nil)
	assert.NoError(t, w.Start(in, out, NewFakeTransformer(), consumeAndCheckErrors(t)))

	go in.PushAll(sampleLogEntries)

	i := 0
	for _ = range out.Messages {
		i++
		if i == FakeInputterMessageCount {
			break
		}
	}
	assert.Equal(t, FakeInputterMessageCount, i)

	assert.NoError(t, w.Close())
}
