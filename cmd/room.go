package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	pb "github.com/hphphp123321/mahjong-common/services/mahjong/v1"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
)

func roomRecv(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			readyReply, err := Client.ReadyStream.Recv()
			if err == io.EOF {
				Client.ReadyDone <- nil
				return
			} else if err != nil {
				Client.ReadyDone <- err
				return
			}
			log.Printf("ReadyStream.Recv: %s", readyReply.Message)
			switch readyReply.GetReply().(type) {
			case *pb.ReadyReply_GetReady:
				Client.HandleGetReadyReply(readyReply)
			}
		}
	}()
	return
}

func roomSend(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer Client.ReadyStream.CloseSend()
		err := selectSend()
		if err != nil {
			Client.ReadyDone <- err
			return
		}
	}()
	return
}

func getOptions() []*survey.Question {
	var allOptions []*survey.Question

	readyOption := &survey.Question{
		Name:   "ready",
		Prompt: &survey.Confirm{Message: "Ready?"},
		Validate: func(val interface{}) error {
			answer := val.(bool)
			if answer {
				Client.P.Ready = true
			} else {
				Client.P.Ready = false
			}
			return nil
		},
	}

	if Client.P.Ready {
		allOptions = append(allOptions, &survey.Question{
			Name:   "cancelReady",
			Prompt: &survey.Confirm{Message: "Cancel Ready?"},
			Validate: func(val interface{}) error {
				answer := val.(bool)
				if answer {
					Client.P.Ready = false
				}
				return nil
			},
		})
	} else {
		allOptions = append(allOptions, readyOption)
	}
	if Client.Room.Owner.PlayerName == Client.P.PlayerName {
		allOptions = append(allOptions, &survey.Question{
			Name: "removePlayer",
			Prompt: &survey.Select{
				Message: "Select player to remove:",
				Options: Client.Room.GetPlayerNames(),
			},
			Validate: func(val interface{}) error {
				playerName := val.(string)
				if playerName == Client.P.PlayerName {
					return fmt.Errorf("You cannot remove yourself from the room")
				}
				// TODO Remove Player
				//Client.Room.RemovePlayer(playerName)
				return nil
			},
		})
		allOptions = append(allOptions, &survey.Question{
			Name: "addRobot",
			Prompt: &survey.Input{
				Message: "Enter robot name:",
			},
			Validate: func(val interface{}) error {
				//robotName := val.(string)
				// TODO Add Robot
				//Client.Room.AddRobot(robotName)
				return nil
			},
		})
		if Client.Room.CheckAllReady() {
			allOptions = append(allOptions, &survey.Question{
				Name:   "startGame",
				Prompt: &survey.Confirm{Message: "Start game?"},
				Validate: func(val interface{}) error {
					answer := val.(bool)
					if answer {
						// TODO start game
					}
					return nil
				},
			})
		}
	}
	allOptions = append(allOptions, &survey.Question{
		Name: "leaveRoom",
		Prompt: &survey.Confirm{
			Message: "Are you sure you want to leave the room?",
			Default: true,
		},
		Validate: func(val interface{}) error {
			answer := val.(bool)
			if answer {
				// TODO Leave Room
				//Client.LeaveRoom()
			}
			return nil
		},
	})
	allOptions = append(allOptions, &survey.Question{
		Name: "chat",
		Prompt: &survey.Input{
			Message: "Enter message:",
		},
		Validate: func(val interface{}) error {
			message := val.(string)
			// TODO Chat
			return Client.SendChatMsg(message)
		},
	})

	return allOptions
}

func selectSend() error {
	for {
		options := getOptions()

		var optionIndex int
		optionSelect := &survey.Select{
			Message: "Select an option:",
			Options: make([]string, len(options)),
		}

		for i, option := range options {
			optionSelect.Options[i] = option.Name
		}

		if err := survey.AskOne(optionSelect, &optionIndex); err != nil {
			return err
		}

		selectedOption := options[optionIndex]
		var answer interface{}

		if err := survey.Ask([]*survey.Question{selectedOption}, &answer); err != nil {
			return err
		}
	}
}
