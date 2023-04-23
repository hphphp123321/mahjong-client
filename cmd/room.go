package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
)

func getOptions() []*survey.Question {
	var allOptions []*survey.Question

	allOptions = append(allOptions, &survey.Question{
		Name: "chat",
		Prompt: &survey.Input{
			Message: "Enter message:",
		},
		Validate: func(val interface{}) error {
			message := val.(string)
			return Client.ReadyChat(message)
		},
	})

	if !Client.IsReady() {
		allOptions = append(allOptions, &survey.Question{
			Name:   "getReady",
			Prompt: &survey.Confirm{Message: "Ready?"},
			Validate: func(val interface{}) error {
				answer := val.(bool)
				if answer {
					return Client.GetReady()
				}
				return nil
			},
		})
	} else {
		allOptions = append(allOptions, &survey.Question{
			Name:   "cancelReady",
			Prompt: &survey.Confirm{Message: "Cancel Ready?"},
			Validate: func(val interface{}) error {
				answer := val.(bool)
				if answer {
					return Client.CancelReady()
				}
				return nil
			},
		})
	}
	if Client.IsOwner() {
		nameMap := Client.Room.GetPlayerNames()
		names := make([]string, 0)
		for name, _ := range nameMap {
			names = append(names, name)
		}
		allOptions = append(allOptions, &survey.Question{
			Name: "removePlayer",
			Prompt: &survey.Select{
				Message: "Select player to remove:",
				Options: names,
			},
			Validate: func(val interface{}) error {
				playerName := val.(string)
				if playerName == Client.Player.Name {
					return fmt.Errorf("you cannot remove yourself from the room")
				}
				return Client.RemovePlayer(nameMap[playerName])
			},
		})
		allOptions = append(allOptions, &survey.Question{
			Name: "addRobot",
			Prompt: &survey.Input{
				Message: "Enter robot name:",
			},
			Validate: func(val interface{}) error {
				robotType := val.(string)
				return Client.AddRobot(robotType, Client.Room.GetIdleSeat())
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
		Name: "listRobots",
		Prompt: &survey.Confirm{
			Message: "Are you sure you want to list robots?",
			Default: true,
		},
		Validate: func(val interface{}) error {
			answer := val.(bool)
			if answer {
				robots, err := Client.ListRobots()
				if err != nil {
					return err
				}
				for _, robot := range robots {
					fmt.Println(robot)
				}
			}
			return nil
		},
	})

	allOptions = append(allOptions, &survey.Question{
		Name: "leaveRoom",
		Prompt: &survey.Confirm{
			Message: "Are you sure you want to leave the room?",
			Default: true,
		},
		Validate: func(val interface{}) error {
			answer := val.(bool)
			if answer {
				return Client.LeaveRoom()
			}
			return nil
		},
	})

	return allOptions
}

func roomSelectSend() error {
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
		switch selectedOption.Prompt.(type) {
		case *survey.Input:
			var inputAnswer string
			answer = &inputAnswer
		case *survey.Select:
			var selectAnswer string
			answer = &selectAnswer
		case *survey.Confirm:
			var confirmAnswer bool
			answer = &confirmAnswer
		default:
			return fmt.Errorf("Unsupported prompt type")
		}

		if err := survey.Ask([]*survey.Question{selectedOption}, answer); err != nil {
			return err
		}
	}
}
