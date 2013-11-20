package wire

import (
  "fmt"
  "github.com/snormore/gologger"
  "github.com/snormore/gowire/message"
  "launchpad.net/tomb"
  "strings"
  "time"
)

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

func (out FakeOutputter) Start(t *tomb.Tomb) error {
  return nil
}

func (out FakeOutputter) Push(msg message.Message) error {
  out.Messages <- msg
  return nil
}
