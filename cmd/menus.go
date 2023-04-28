package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/hphphp123321/mahjong-client/app/errs"
	"github.com/hphphp123321/mahjong-client/app/service/client"
	"github.com/hphphp123321/mahjong-go/mahjong"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
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
	readyStream, err := Client.Client.Ready(Client.Ctx)
	if err != nil {
		log.Warnln("start ready stream error:", err)
		Client.Room = nil
		if err := Client.Player.LeaveRoom(); err != nil {
			log.Warnln("leave room error:", err)
		}
		return
	}
	Client.ReadyStream = readyStream
	log.Debugln("start stream send")
	wg := sync.WaitGroup{}
	done := client.StartReadyRecvStream(Client, &wg)
	go RefreshRoom(Client)
	for {
		err := roomSelectSend(done, &wg)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorln(err)
		}
	}
	Client.ReadyStream = nil
	if err == errs.ErrGameStart {
		gameMenu()
	}
}

func gameMenu() {
	gameStream, err := Client.Client.Game(Client.Ctx)
	if err != nil {
		log.Warnln("start game stream error:", err)
		return
	}
	Client.GameStream = gameStream
	log.Debugln("start stream recv")
	actionChan := make(chan mahjong.Calls)
	done := client.StartGameRecvStream(Client, actionChan)
	go RefreshGame(Client)
	for {
		err := gameSelectSend(done, actionChan)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorln(err)
		}
	}
}
