package cmd

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"time"
)

func chooseAction(actions mahjong.Calls) (mahjong.Call, error) {
	var action mahjong.Call
	var actionOptions []string
	for _, call := range actions {
		actionOptions = append(actionOptions, call.String())
	}
	return action, nil
}

func gameSelectSend(done chan error, actionChan chan mahjong.Calls) error {
	var actions mahjong.Calls
	select {
	case err := <-done:
		return err
	case <-time.After(Client.Delay*10 + time.Millisecond*100):
	case actions = <-actionChan:

	}
}
