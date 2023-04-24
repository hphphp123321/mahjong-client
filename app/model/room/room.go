package room

import (
	"github.com/hphphp123321/mahjong-client/app/errs"
	"github.com/hphphp123321/mahjong-client/app/model/player"
)

type Room struct {
	ID        string
	Name      string
	Players   map[int]*player.Player
	OwnerSeat int
}

func NewRoom(p *player.Player, name string, id string) (*Room, error) {
	r := &Room{
		ID:        id,
		Name:      name,
		Players:   map[int]*player.Player{1: p},
		OwnerSeat: 1,
	}
	if err := p.JoinRoom(r.ID, 1); err != nil {
		return nil, err
	}
	return r, nil
}

func NewRoomFromInfo(roomInfo *Info) *Room {
	players := map[int]*player.Player{}
	for _, info := range roomInfo.PlayerInfos {
		p := player.NewPlayerByInfo(info)
		p.RoomID = roomInfo.ID
		players[info.Seat] = p
	}
	return &Room{
		ID:        roomInfo.ID,
		Name:      roomInfo.Name,
		Players:   players,
		OwnerSeat: roomInfo.OwnerSeat,
	}
}

func (r *Room) Refresh(info *Info) error {
	if r.ID != info.ID {
		return errs.ErrRoomIDNotMatch
	}
	if r.Name != info.Name {
		return errs.ErrRoomNameNotMatch
	}
	if r.OwnerSeat != info.OwnerSeat {
		return errs.ErrRoomOwnerSeatNotMatch
	}
	for _, info := range info.PlayerInfos {
		if p, ok := r.Players[info.Seat]; ok {
			if err := p.Refresh(info); err != nil {
				return err
			}
		} else {
			return errs.ErrPlayerNotInRoom
		}
	}
	return nil
}

func (r *Room) Join(p *player.Player, seat int) error {
	if len(r.Players) == 4 {
		return errs.ErrRoomFull
	}
	if err := p.JoinRoom(r.ID, seat); err != nil {
		return err
	}
	r.Players[seat] = p
	return nil
}

func (r *Room) Leave(seat int) error {
	if p, ok := r.Players[seat]; ok {
		if err := p.LeaveRoom(); err != nil {
			return err
		}
		delete(r.Players, seat)
		return nil
	}
	return errs.ErrPlayerNotInRoom
}

func (r *Room) SetReady(seat int) error {
	if p, ok := r.Players[seat]; ok {
		p.Ready = true
		return nil
	}
	return errs.ErrPlayerNotInRoom
}

func (r *Room) SetUnReady(seat int) error {
	if p, ok := r.Players[seat]; ok {
		p.Ready = false
		return nil
	}
	return errs.ErrPlayerNotInRoom
}

func (r *Room) GetPlayerNames() map[string]int {
	var nameMap = make(map[string]int)
	for _, p := range r.Players {
		nameMap[p.Name] = p.Seat
	}
	return nameMap
}

func (r *Room) GetIdleSeat() int {
	for _, seat := range []int{1, 2, 3, 4} {
		if _, ok := r.Players[seat]; !ok {
			return seat
		}
	}
	panic("room is full")
}

func (r *Room) CheckAllReady() bool {
	if len(r.Players) < 4 {
		return false
	}
	for _, p := range r.Players {
		if !p.Ready {
			return false
		}
	}
	return true
}
