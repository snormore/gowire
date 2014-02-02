package wire

import (
	"launchpad.net/tomb"
)

type FakeInputter struct {
	Messages chan interface{}
}

func NewFakeInputter() *FakeInputter {
	in := new(FakeInputter)
	in.Messages = make(chan interface{}, 1024)
	return in
}

func (in *FakeInputter) PushAll(messages []string) error {
	for _, msg := range messages {
		in.Messages <- msg
	}
	return nil
}

func (in *FakeInputter) Transform(rawMessage interface{}) (interface{}, error) {
	return rawMessage, nil
}

func (in *FakeInputter) Listen() chan interface{} {
	return in.Messages
}

func (in *FakeInputter) Start(t *tomb.Tomb) error {
	return nil
}

func (in *FakeInputter) FinalizeMessage(msg interface{}) error {
	return nil
}

func (in *FakeInputter) Close() error {
	return nil
}

type FakeOutputter struct {
	Messages chan interface{}
}

func (out *FakeOutputter) Start(t *tomb.Tomb) error {
	return nil
}

func (out *FakeOutputter) Push(msg interface{}) error {
	out.Messages <- msg
	return nil
}

func (out *FakeOutputter) Close() error {
	return nil
}
