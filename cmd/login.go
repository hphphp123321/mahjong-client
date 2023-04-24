package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
)

func login() error {
	usernamePrompt := &survey.Input{
		Message: "Enter your username:",
	}
	var username string
	if err := survey.AskOne(usernamePrompt, &username); err != nil {
		return err
	}

	log.Println("Logging in as %s...\n", username)
	if err := Client.Login(username); err != nil {
		return err
	}
	go Ping(Client)
	SetupSignalHandler(Client)
	return nil
}

func quit() error {
	log.Println("Goodbye!")
	if Client != nil {
		if err := Client.Logout(); err != nil {
			return err
		}
	}
	return nil
}
