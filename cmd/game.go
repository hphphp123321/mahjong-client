package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/hphphp123321/mahjong-go/mahjong"
	log "github.com/sirupsen/logrus"
	"time"
)

func callOption(call *mahjong.Call) string {
	option := ""
	callType := call.CallType
	option += callType.String() + ": "
	switch callType {
	case mahjong.Discard:
		option += call.CallTiles[0].UTF8()
	case mahjong.Pon:
		tiles := call.CallTiles[:3]
		option += tiles.String() + " "
	case mahjong.Chi:
		tiles := call.CallTiles[:3]
		option += tiles.String() + " "
	case mahjong.DaiMinKan:
		tiles := call.CallTiles[:4]
		option += tiles.String() + " "
	case mahjong.ShouMinKan:
		tiles := call.CallTiles[:4]
		option += tiles.String() + " "
	case mahjong.AnKan:
		tiles := call.CallTiles[:4]
		option += tiles.String() + " "
	case mahjong.Riichi:
		option += call.CallTiles[0].String() + " "
	case mahjong.Ron | mahjong.Tsumo | mahjong.ChanKan:
		option += call.CallTiles[0].String() + " "
	case mahjong.KyuuShuKyuuHai | mahjong.Next | mahjong.Skip:
		break
	}
	return option
}

func chooseAction(actions mahjong.Calls) *mahjong.Call {
	var action *mahjong.Call
	var actionOptions []string
	for _, call := range actions {
		actionOptions = append(actionOptions, callOption(call))
	}
	var optionIdx int
	optionSelect := &survey.Select{
		Message:  "Choose an action:",
		Options:  actionOptions,
		PageSize: 20,
	}
	err := survey.AskOne(optionSelect, &optionIdx)
	if err != nil {
		log.Warnln("Choose action error:", err)
	}
	action = actions[optionIdx]
	return action
}

func getState(state *mahjong.BoardState) string {
	s := "hand tiles: " + state.HandTiles.String() + "\n"
	s += "melds: "
	for _, meld := range state.PlayerStates[state.PlayerWind].Melds {
		s += meld.CallTiles.String() + "; "
	}
	return s
}

func gameSelectSend(done chan error, actionChan chan mahjong.Calls) error {
	var actions mahjong.Calls
	select {
	case err := <-done:
		return err
	case <-time.After(Client.Delay*10 + time.Millisecond*1000):
	case actions = <-actionChan:
		state := getState(Client.BoardState)
		log.Infoln(state)
		action := chooseAction(actions)
		if action != nil {
			if err := Client.SendAction(action); err != nil {
				log.Warnln("Send action error:", err)
			}
		} else {
			log.Warnln("No action selected")
		}
	}
	return nil
}