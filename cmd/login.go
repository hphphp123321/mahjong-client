package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	"github.com/hphphp123321/mahjong-common/player"
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
	Client.P = player.NewPlayer(username, uuid.Nil)
	if err := Client.Login(); err != nil {
		return err
	}
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
