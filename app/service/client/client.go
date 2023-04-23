package client

import (
	"context"
	"github.com/hphphp123321/go-common"
	pb "github.com/hphphp123321/mahjong-client/app/api/v1"
	"github.com/hphphp123321/mahjong-client/app/errs"
	"github.com/hphphp123321/mahjong-client/app/model/player"
	"github.com/hphphp123321/mahjong-client/app/model/room"
	log "github.com/sirupsen/logrus"
	"time"
)

type MahjongClient struct {
	ID  string
	Ctx context.Context

	Client      pb.MahjongClient
	ReadyStream pb.Mahjong_ReadyClient
	GameStream  pb.Mahjong_GameClient

	Player *player.Player

	Room      *room.Room
	RoomInfos []*room.Info

	Delay time.Duration
}

func NewMahjongClient(id string) *MahjongClient {
	return &MahjongClient{
		ID: id,
	}
}

func (c *MahjongClient) Ping() error {
	start := time.Now()
	if _, err := c.Client.Ping(c.Ctx, &pb.Empty{}); err != nil {
		return err
	}
	end := time.Now()
	usedTime := end.Sub(start)
	c.Delay = usedTime
	return nil
}

func (c *MahjongClient) Login(name string) error {
	reply, err := c.Client.Login(c.Ctx, &pb.LoginRequest{
		PlayerName: name,
	})
	if err != nil {
		return err
	}
	id := reply.GetId()
	if id == "" {
		log.Errorf("Login failed: %s", reply.GetMessage())
	}
	c.ID = id
	c.Player = player.NewPlayer(name)
	return nil
}

func (c *MahjongClient) Logout() error {
	if c.ID == "" {
		return nil
	}
	if c.Room != nil {
		// TODO leave room
	}
	_, err := c.Client.Logout(c.Ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	c.ID = ""
	return nil
}

func (c *MahjongClient) ListRooms(nameFilter string) (map[string]string, error) {
	reply, err := c.Client.ListRooms(c.Ctx, &pb.ListRoomsRequest{
		RoomName: &nameFilter,
	})
	if err != nil {
		return nil, err
	}
	roomInfos := reply.GetRooms()
	c.RoomInfos = common.MapSlice(roomInfos, ToRoomInfo)
	var roomIDNames = make(map[string]string)
	for _, roomInfo := range c.RoomInfos {
		roomIDNames[roomInfo.ID] = roomInfo.Name
	}
	return roomIDNames, nil
}

func (c *MahjongClient) CreateRoom(name string) error {
	reply, err := c.Client.CreateRoom(c.Ctx, &pb.CreateRoomRequest{
		RoomName: name,
	})
	if err != nil {
		return err
	}
	roomInfo := ToRoomInfo(reply.GetRoom())
	c.Room = room.NewRoomFromInfo(roomInfo)
	c.ReadyStream, err = c.Client.Ready(c.Ctx)
	if err != nil {
		return err
	}
	log.Debugln("start stream send")
	return nil
}

func (c *MahjongClient) JoinRoom(id string) error {
	reply, err := c.Client.JoinRoom(c.Ctx, &pb.JoinRoomRequest{
		RoomID: id,
	})
	if err != nil {
		return err
	}
	roomInfo := ToRoomInfo(reply.GetRoom())
	c.Room = room.NewRoomFromInfo(roomInfo)
	if err := c.Player.JoinRoom(c.Room.ID, c.Room.OwnerSeat); err != nil {
		return err
	}
	c.ReadyStream, err = c.Client.Ready(c.Ctx)
	if err != nil {
		return err
	}
	log.Debugln("start stream send")
	return nil
}

func (c *MahjongClient) GetReady() error {
	if c.Player == nil {
		return errs.ErrPlayerNotFound
	}
	if c.Player.Ready {
		return errs.ErrPlayerReady
	}
	if err := c.ReadyStream.Send(&pb.ReadyRequest{
		Request: &pb.ReadyRequest_GetReady{
			GetReady: &pb.Empty{},
		},
	}); err != nil {
		return err
	}
	// TODO c.Player.Ready = true
	return nil
}

func (c *MahjongClient) CancelReady() error {
	if c.Player == nil {
		return errs.ErrPlayerNotFound
	}
	if !c.Player.Ready {
		return errs.ErrPlayerNotReady
	}
	if err := c.ReadyStream.Send(&pb.ReadyRequest{
		Request: &pb.ReadyRequest_CancelReady{
			CancelReady: &pb.Empty{},
		},
	}); err != nil {
		return err
	}
	// TODO c.Player.Ready = false
	return nil
}

func (c *MahjongClient) AddRobot(robotType string, robotSeat int) error {
	if c.Room == nil {
		return errs.ErrRoomNotFound
	}
	if c.Room.OwnerSeat != c.Player.Seat {
		return errs.ErrPlayerNotOwner
	}
	if err := c.ReadyStream.Send(&pb.ReadyRequest{
		Request: &pb.ReadyRequest_AddRobot{
			AddRobot: &pb.AddRobotRequest{
				RobotType: robotType,
				RobotSeat: int32(robotSeat),
			},
		},
	}); err != nil {
		return err
	}
	// TODO if err := c.Room.Join(player.NewPlayer(robotType), robotSeat); err != nil {
	//	return err
	//}
	return nil
}

func (c *MahjongClient) RemovePlayer(seat int) error {
	if c.Room == nil {
		return errs.ErrRoomNotFound
	}
	if c.Room.OwnerSeat != c.Player.Seat {
		return errs.ErrPlayerNotOwner
	}
	if err := c.ReadyStream.Send(&pb.ReadyRequest{
		Request: &pb.ReadyRequest_RemovePlayer{
			RemovePlayer: &pb.RemovePlayerRequest{
				PlayerSeat: int32(seat),
			},
		},
	}); err != nil {
		return err
	}
	// TODO if err := c.Room.Leave(seat); err != nil {
	//	return err
	//}
	return nil
}

func (c *MahjongClient) LeaveRoom() error {
	if c.Room == nil {
		return errs.ErrRoomNotFound
	}
	if err := c.ReadyStream.Send(&pb.ReadyRequest{
		Request: &pb.ReadyRequest_LeaveRoom{
			LeaveRoom: &pb.Empty{},
		},
	}); err != nil {
		return err
	}
	// TODO if err := c.Player.LeaveRoom(); err != nil {
	//	return err
	//}
	//c.Room = nil
	return nil
}

func (c *MahjongClient) ReadyChat(msg string) error {
	if c.Room == nil {
		return errs.ErrRoomNotFound
	}
	if err := c.ReadyStream.Send(&pb.ReadyRequest{
		Request: &pb.ReadyRequest_Chat{
			Chat: &pb.ChatRequest{
				Message: msg,
			},
		},
	}); err != nil {
		return err
	}
	return nil
}

func (c *MahjongClient) ListRobots() ([]string, error) {
	if c.Room == nil {
		return nil, errs.ErrRoomNotFound
	}
	if c.Room.OwnerSeat != c.Player.Seat {
		return nil, errs.ErrPlayerNotOwner
	}
	reply, err := c.Client.ListRobots(c.Ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	return reply.GetRobotTypes(), nil
}