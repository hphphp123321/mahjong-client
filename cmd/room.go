package cmd

import (
	"github.com/AlecAivazis/survey/v2"
)

func getOptions() []string {
	var allOptions []string
	allOptions = append(allOptions, "chat")
	if !Client.IsReady() {
		allOptions = append(allOptions, "get ready")
	} else {
		allOptions = append(allOptions, "cancel ready")
	}
	if Client.IsOwner() {
		allOptions = append(allOptions, "remove player")
		allOptions = append(allOptions, "add robot")
		if Client.Room.CheckAllReady() {
			allOptions = append(allOptions, "start game")
		}
	}
	allOptions = append(allOptions, "leave room")
	return allOptions
}

func roomSelectSend() error {
	options := getOptions()

	var optionIndex int
	optionSelect := &survey.Select{
		Message: "Select an option:",
		Options: options,
	}

	if err := survey.AskOne(optionSelect, &optionIndex); err != nil {
		return err
	}

	return nil
}
