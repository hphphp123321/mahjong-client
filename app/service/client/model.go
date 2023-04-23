package client

import (
	"github.com/hphphp123321/go-common"
	pb "github.com/hphphp123321/mahjong-client/app/api/v1"
	"github.com/hphphp123321/mahjong-client/app/model/player"
	"github.com/hphphp123321/mahjong-client/app/model/room"
)

func ToPlayerInfo(info *pb.PlayerInfo) *player.Info {
	return &player.Info{
		Name:  info.PlayerName,
		Seat:  int(info.PlayerSeat),
		Ready: info.IsReady,
	}
}

func ToRoomInfo(info *pb.RoomInfo) *room.Info {
	return &room.Info{
		ID:          info.RoomID,
		Name:        info.RoomName,
		OwnerSeat:   int(info.OwnerSeat),
		PlayerInfos: common.MapSlice(info.Players, ToPlayerInfo),
	}
}
