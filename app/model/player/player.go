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

func NewRobot(robotType string) *Player {
	return &Player{
		Name:  robotType,
		Ready: true,
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

func (p *Player) Refresh(info *Info) error {
	if p.Name != info.Name {
		return errs.ErrPlayerNameNotMatch
	}
	if p.Seat != info.Seat {
		return errs.ErrPlayerSeatNotMatch
	}
	p.Ready = info.Ready
	return nil

}
