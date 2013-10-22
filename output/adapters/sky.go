package sky_ouputter

import (
	"github.com/snormore/gowire"
	"github.com/snormore/gowire/message"
	"sync"
)

type SkyOutputter struct{}

func (l SkyOutputter) Start(messages chan message.Message, wg *sync.WaitGroup, t *tomb.Tomb) error {
	return nil
}
