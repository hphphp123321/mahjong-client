package player

import "github.com/hphphp123321/mahjong-client/app/errs"

type Player struct {
	Name   string
	Seat   int
	RoomID string
	Ready  bool
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
	}
}

func NewPlayerByInfo(info *Info) *Player {
	return &Player{
		Name:  info.Name,
		Seat:  info.Seat,
		Ready: info.Ready,
	}
}

func (p *Player) JoinRoom(roomID string, seat int) error {
	if p.RoomID != "" {
		return errs.ErrPlayerAlreadyInRoom
	}
	p.RoomID = roomID
	p.Seat = seat
	return nil
}

func (p *Player) LeaveRoom() error {
	if p.RoomID == "" {
		return errs.ErrPlayerNotInRoom
	}
	p.RoomID = ""
	p.Seat = 0
	return nil
}
