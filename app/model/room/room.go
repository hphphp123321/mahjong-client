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
		players[info.Seat] = player.NewPlayerByInfo(info)
	}
	return &Room{
		ID:        roomInfo.ID,
		Name:      roomInfo.Name,
		Players:   players,
		OwnerSeat: roomInfo.OwnerSeat,
	}
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
