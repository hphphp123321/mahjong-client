package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"time"
)

var (
	doNothingOption    = "do nothing"
	chatOption         = "chat"
	getReadyOption     = "get ready"
	cancelReadyOption  = "cancel ready"
	removePlayerOption = "remove player"
	addRobotOption     = "add robot"
	startGameOption    = "start game"
	leaveRoomOption    = "leave room"
)

func getOptions() []string {
	var allOptions []string
	allOptions = append(allOptions, doNothingOption)
	allOptions = append(allOptions, chatOption)
	if !Client.IsReady() {
		allOptions = append(allOptions, getReadyOption)
		allOptions = append(allOptions, cancelReadyOption)
	} else {
		allOptions = append(allOptions, cancelReadyOption)
	}
	if Client.IsOwner() {
		allOptions = append(allOptions, removePlayerOption)
		if !Client.Room.IsFull() {
			allOptions = append(allOptions, addRobotOption)
		}
		if Client.Room.CheckAllReady() {
			allOptions = append(allOptions, startGameOption)
		}
	}
	return allOptions
}

func chat() error {
	var message string
	prompt := &survey.Input{
		Message: "Enter your message:",
	}
	if err := survey.AskOne(prompt, &message); err != nil {
		return err
	}
	if err := Client.ReadyChat(message); err != nil {
		return err
	}
	return nil
}

func getReady() error {
	var confirm bool
	prompt := &survey.Confirm{
		Message: "Are you sure to get ready?",
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return err
	}
	if confirm {
		if err := Client.GetReady(); err != nil {
			return err
		}
	}
	return nil
}

func cancelReady() error {
	var confirm bool
	prompt := &survey.Confirm{
		Message: "Are you sure to cancel ready?",
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return err
	}
	if confirm {
		if err := Client.CancelReady(); err != nil {
			return err
		}
	}
	return nil
}

func removePlayer() error {
	var options = make([]string, 0)
	names := Client.Room.GetPlayerNames()
	for name, seat := range names {
		if seat == Client.Player.Seat {
			continue
		}
		options = append(options, name)
	}
	options = append(options, "cancel")
	var optionName string
	prompt := &survey.Select{
		Message: "Select a player to remove:",
		Options: options,
	}

	if err := survey.AskOne(prompt, &optionName); err != nil {
		return err
	}
	if optionName == "cancel" {
		return nil
	}

	return Client.RemovePlayer(names[optionName])
}

func addRobot() error {
	robots, err := Client.ListRobots()
	if err != nil {
		return err
	}
	var options = make([]string, 0)
	for _, robot := range robots {
		options = append(options, robot)
	}
	options = append(options, "cancel")
	var optionName string
	prompt := &survey.Select{
		Message: "Select a robot to add:",
		Options: options,
	}
	if err := survey.AskOne(prompt, &optionName); err != nil {
		return err
	}
	if optionName == "cancel" {
		return nil
	}
	seat := Client.Room.GetIdleSeat()
	return Client.AddRobot(optionName, seat)
}

func startGame() error {
	var confirm bool
	prompt := &survey.Confirm{
		Message: "Are you sure to start game?",
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return err
	}
	if confirm {
		if err := Client.StartGame(); err != nil {
			return err
		}
	}
	return nil
}

func leaveRoom() error {
	var confirm bool
	prompt := &survey.Confirm{
		Message: "Are you sure to leave room?",
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return err
	}
	if confirm {
		if err := Client.LeaveRoom(); err != nil {
			return err
		}
	}
	return nil
}

func chooseRoomAction(done chan error) (errChan chan error) {
	errChan = make(chan error)
	go func() {
		for {
			options := getOptions()

			var optionName string
			optionSelect := &survey.Select{
				Message: "Select an option:",
				Options: options,
			}
			select {
			case <-done:
				return
			default:
				if err := survey.AskOne(optionSelect, &optionName); err != nil {
					done <- err
					return
				}
				switch optionName {
				case doNothingOption:
					time.Sleep(time.Millisecond*500 + Client.Delay*10)
					continue
				case chatOption:
					if err := chat(); err != nil {
						errChan <- err
						return
					}
				case getReadyOption:
					if err := getReady(); err != nil {
						errChan <- err
						return
					}
					time.Sleep(time.Millisecond*500 + Client.Delay*10)
				case cancelReadyOption:
					if err := cancelReady(); err != nil {
						errChan <- err
						return
					}
					time.Sleep(time.Millisecond*500 + Client.Delay*10)
				case removePlayerOption:
					if err := removePlayer(); err != nil {
						errChan <- err
						return
					}
					time.Sleep(time.Millisecond*500 + Client.Delay*10)
				case addRobotOption:
					if err := addRobot(); err != nil {
						errChan <- err
						return
					}
					time.Sleep(time.Millisecond*500 + Client.Delay*10)
				case startGameOption:
					if err := startGame(); err != nil {
						errChan <- err
						return
					}
					time.Sleep(time.Second * 1)
				case leaveRoomOption:
					if err := leaveRoom(); err != nil {
						errChan <- err
						return
					}
					time.Sleep(time.Millisecond*500 + Client.Delay*10)
				}
			}
		}
	}()
	return errChan
}

func roomSelectSend(done chan error) error {
	optionDone := make(chan error)
	errChan := chooseRoomAction(optionDone)

	for {
		select {
		case err := <-done:
			optionDone <- err
			return err
		case err := <-errChan:
			done <- err
			return err
		case <-time.After(Client.Delay*10 + time.Millisecond*1000):
		}
	}
}
