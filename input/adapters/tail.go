package input_adapter

import (
	"encoding/json"
	"github.com/snormore/gologger"
	"github.com/snormore/gotail"
	"github.com/snormore/gowire/config"
	"github.com/snormore/gowire/message"
	"io/ioutil"
	"launchpad.net/tomb"
	"math"
	"os"
	"strings"
)

const (
	LatestFinalizedFilePath = "./tmp/local.message"
	EmptyEventId            = "-"
)

type TailInputter struct {
	config.Configurable

	tailer    *tail.Tailer
	lines     chan interface{}
	finalized chan message.Message
}

type TailConfig struct {
	FilePath                string
	StartEvent              string
	LatestFinalizedFlushMod int
}

var DefaultTailConfig = TailConfig{
	FilePath:                "logs/development.json.log",
	StartEvent:              "",
	LatestFinalizedFlushMod: 100,
}

// func createTailConfig(c TailConfig) TailConfig {
// 	// merge given config and defaults
// }

type Event struct {
	Id string `json:"event_id"`
}

func NewTailInputter(conf interface{}) *TailInputter {
	in := TailInputter{}
	in.Config = conf.(config.Config)
	in.tailer = tail.NewTailer()
	in.lines = make(chan interface{}, tail.LinesChannelSize)
	go func() {
		for line := range in.tailer.Listen() {
			in.lines <- line
		}
	}()
	return &in
}

func (in TailInputter) Start(t *tomb.Tomb) error {
	return in.tailer.Read(in.Config.(TailConfig).FilePath, in.Config.(TailConfig).StartEvent, t)
}

func (in TailInputter) Transform(rawMessage interface{}) (message.Message, error) {
	event := new(Event)
	if err := json.Unmarshal([]byte(rawMessage.(string)), event); err != nil {
		return message.Message{"", rawMessage}, err
	}
	return message.Message{event.Id, rawMessage}, nil
}

func (in TailInputter) Listen() chan interface{} {
	return in.lines
}

func (in TailInputter) FinalizeMessage(msg message.Message) error {
	in.finalized <- msg
	return nil
}

func (in *TailInputter) getLatestFinalizedMessageId() (string, error) {
	if _, err := os.Stat(LatestFinalizedFilePath); os.IsNotExist(err) {
		logger.Info("No such file or directory: %s", LatestFinalizedFilePath)
		return "", nil
	}
	eventBytes, err := ioutil.ReadFile(LatestFinalizedFilePath)
	if err != nil {
		logger.Panic("Failed to open latest finalized message file: %+v", err)
	}
	event := string(eventBytes)
	event = strings.Trim(string(event), " \n\r\t")
	if string(event) == EmptyEventId {
		event = ""
	}
	return "", nil
}

func (in *TailInputter) setLatestFinalizedMessageId(msgId string) error {
	return ioutil.WriteFile(LatestFinalizedFilePath, []byte(msgId), 0644)
}

func (in *TailInputter) finalizeListener(t *tomb.Tomb) {
	msgCounter := 0
	lastFinalizedMsgId := ""
	for {
		select {
		case msg := <-in.finalized:
			if math.Mod(float64(msgCounter), float64(in.Config.(TailConfig).LatestFinalizedFlushMod)) == 0.0 && strings.Trim(msg.Id, " \n\r\t") != "" {
				logger.VerboseDebug("Saving latest finalized message ID %s...", msg.Id)
				in.setLatestFinalizedMessageId(msg.Id)
			}
			if msg.Id != "" {
				lastFinalizedMsgId = msg.Id
			}
			msgCounter++
		case <-t.Dying():
			for msg := range in.finalized {
				if msg.Id != "" {
					lastFinalizedMsgId = msg.Id
				}
			}
			if lastFinalizedMsgId != "" {
				logger.Info("Saving latest finalized message ID %s before exit...", lastFinalizedMsgId)
				in.setLatestFinalizedMessageId(lastFinalizedMsgId)
			}
			return
		}
	}
}
