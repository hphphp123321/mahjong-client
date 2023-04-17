package cmd

import (
	"errors"
	"fmt"
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
	if err := Client.RefreshRoom(roomName); err != nil {
		return err
	}
	for _, room := range Client.RoomList {
		fmt.Println(room.RoomName)
	}
	return nil
}

func joinRoom() error {
	var options = make([]string, 0)
	for _, room := range Client.RoomList {
		options = append(options, room.RoomName)
	}
	if len(options) == 0 {
		return errors.New("no room found")
	}
	roomNumberPrompt := &survey.Select{
		Message: "Enter the room name:",
		Options: options,
	}
	var roomName string
	if err := survey.AskOne(roomNumberPrompt, &roomName); err != nil {
		return err
	}
	for _, room := range Client.RoomList {
		if room.RoomName == roomName {
			return Client.JoinRoom(room.RoomID.String())
		}
	}
	return errors.New("room not found")
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
	Client = nil
	// 在此处添加注销逻辑
}
