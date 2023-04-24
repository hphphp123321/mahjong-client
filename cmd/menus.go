package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/hphphp123321/mahjong-client/app/service/client"
	log "github.com/sirupsen/logrus"
)

func logMenu() {
	for {
		menu := &survey.Select{
			Message: "Please choose an option:",
			Options: []string{"Login", "Quit"},
		}
		var choice string
		err := survey.AskOne(menu, &choice)
		if err != nil {
			return
		}

		switch choice {
		case "Login":
			if err := login(); err != nil {
				log.Fatal(err)
			}
			lobbyMenu()
		case "Quit":
			if err := quit(); err != nil {
				log.Fatal(err)
			}
			return
		}
	}
}

func lobbyMenu() {
	// ... 省略 ...

	for {
		prompt := &survey.Select{
			Message: "Please choose an option:",
			Options: []string{"ListRooms", "JoinRoom", "CreateRoom", "Logout"},
		}
		var action string
		err := survey.AskOne(prompt, &action)
		if err != nil {
			return
		}
		_, err = Client.ListRooms("")
		if err != nil {
			log.Errorln(err)
			return
		}
		switch action {
		case "ListRooms":
			err := listRooms()
			if err != nil {
				fmt.Println(err)
			}
		case "JoinRoom":
			err := joinRoom()
			if err != nil {
				fmt.Println(err)
				continue
			}
			roomMenu()
		case "CreateRoom":
			err := createRoom()
			if err != nil {
				fmt.Println(err)
				continue
			}
			roomMenu()
		case "Logout":
			logout()
			return
		}
	}
}

func roomMenu() {
	done := client.StartReadyRecvStream(Client)
	go RefreshRoom(Client)
	for {
		select {
		case err := <-done:
			if err != nil {
				log.Errorln(err)
			} else {
				return
			}
		default:
			if err := roomSelectSend(); err != nil {
				log.Errorln(err)
			}
		}
	}
}
