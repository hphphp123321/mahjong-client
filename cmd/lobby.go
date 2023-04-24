package cmd

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
)

func listRooms() error {
	roomNamePrompt := &survey.Input{
		Message: "Enter the room name to find, default is find all:",
	}
	var roomName string
	if err := survey.AskOne(roomNamePrompt, &roomName); err != nil {
		return err
	}
	roomInfos, err := Client.ListRooms(roomName)
	if err != nil {
		return err
	}
	for id, name := range roomInfos {
		log.Infoln("Room ID:", id, "Name: ", name)
	}
	return nil
}

func joinRoom() error {
	var options = make([]string, 0)
	for _, roomInfo := range Client.RoomInfos {
		options = append(options, roomInfo.Name)
	}
	if len(options) == 0 {
		return errors.New("no room found")
	}
	var optionIdx int
	roomNumberPrompt := &survey.Select{
		Message: "Enter the room name:",
		Options: options,
	}

	if err := survey.AskOne(roomNumberPrompt, &optionIdx); err != nil {
		return err
	}
	roomId := Client.RoomInfos[optionIdx].ID
	return Client.JoinRoom(roomId)
}

func createRoom() error {
	roomNamePrompt := &survey.Input{
		Message: "Enter the room name to create:",
	}
	var roomName string
	err := survey.AskOne(roomNamePrompt, &roomName)
	if err != nil {
		return err
	}
	return Client.CreateRoom(roomName)
}

func logout() {
	log.Println("Logging out...")
	if err := Client.Logout(); err != nil {
		log.Fatal(err)
	}
	// 在此处添加注销逻辑
}
