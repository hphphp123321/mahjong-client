package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"sync"
	"time"
)

var (
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
	allOptions = append(allOptions, "chat")
	if !Client.IsReady() {
		allOptions = append(allOptions, "get ready")
		allOptions = append(allOptions, "leave room")
	} else {
		allOptions = append(allOptions, "cancel ready")
	}
	if Client.IsOwner() {
		allOptions = append(allOptions, "remove player")
		if !Client.Room.IsFull() {
			allOptions = append(allOptions, "add robot")
		}
		if Client.Room.CheckAllReady() {
			allOptions = append(allOptions, "start game")
		}
	}
	return allOptions
}

func chat(wg *sync.WaitGroup) error {
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

func getReady(wg *sync.WaitGroup) error {
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
	} else {
		wg.Done()
	}
	return nil
}

func cancelReady(wg *sync.WaitGroup) error {
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
	} else {
		wg.Done()
	}
	return nil
}

func removePlayer(wg *sync.WaitGroup) error {
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
		wg.Done()
		return nil
	}
	wg.Done()
	return Client.RemovePlayer(names[optionName])
}

func addRobot(wg *sync.WaitGroup) error {
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
		wg.Done()
		return nil
	}
	seat := Client.Room.GetIdleSeat()
	return Client.AddRobot(optionName, seat)
}

func startGame(wg *sync.WaitGroup) error {
	var confirm bool
	prompt := &survey.Confirm{
		Message: "Are you sure to start game?",
	}
	if err := survey.AskOne(prompt, &confirm); err != nil {
		return err
	}
	if confirm {
		// TODO: start game
		//if err := Client.StartGame(); err != nil {
		//	return err
		//}
	} else {
		wg.Done()
	}
	return nil
}

func leaveRoom(wg *sync.WaitGroup) error {
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
	} else {
		wg.Done()
	}
	return nil
}

func roomSelectSend(done chan error, wg *sync.WaitGroup) error {
	wg.Wait()

	options := getOptions()

	var optionName string
	optionSelect := &survey.Select{
		Message: "Select an option:",
		Options: options,
	}

	select {
	case err := <-done:
		return err
	case <-time.After(Client.Delay*10 + time.Millisecond*100):
	}

	if err := survey.AskOne(optionSelect, &optionName); err != nil {
		return err
	}

	wg.Add(1)
	switch optionName {
	case chatOption:
		if err := chat(wg); err != nil {
			return err
		}
	case getReadyOption:
		if err := getReady(wg); err != nil {
			return err
		}
	case cancelReadyOption:
		if err := cancelReady(wg); err != nil {
			return err
		}
	case removePlayerOption:
		if err := removePlayer(wg); err != nil {
			return err
		}
	case addRobotOption:
		if err := addRobot(wg); err != nil {
			return err
		}
	case startGameOption:
		if err := startGame(wg); err != nil {
			return err
		}
	case leaveRoomOption:
		if err := leaveRoom(wg); err != nil {
			return err
		}
	}
	return nil
}
