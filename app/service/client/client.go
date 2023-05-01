package client

import "C"
import (
	"context"
	"github.com/hphphp123321/go-common"
	pb "github.com/hphphp123321/mahjong-client/app/api/v1"
	"github.com/hphphp123321/mahjong-client/app/errs"
	"github.com/hphphp123321/mahjong-client/app/model/player"
	"github.com/hphphp123321/mahjong-client/app/model/room"
	"github.com/hphphp123321/mahjong-go/mahjong"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"math/rand"
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

	BoardState *mahjong.BoardState

	Delay time.Duration
}

func NewMahjongClient(ctx context.Context, client pb.MahjongClient) *MahjongClient {
	return &MahjongClient{
		Ctx:    ctx,
		Client: client,
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
	header := metadata.New(map[string]string{"id": id})
	c.Ctx = metadata.NewOutgoingContext(c.Ctx, header)
	c.ID = id
	c.Player = player.NewPlayer(name)
	return nil
}

func (c *MahjongClient) Logout() error {
	if c.ID == "" {
		return nil
	}
	if c.Client == nil {
		return nil
	}
	if c.Room != nil {
		if err := c.LeaveRoom(); err != nil {
			log.Warnln(err)
		}
	}
	time.Sleep(time.Millisecond * 100)
	if c.ReadyStream != nil {
		if err := c.ReadyStream.CloseSend(); err != nil {
			log.Warnln(err)
		}
		c.ReadyStream = nil
	}
	time.Sleep(time.Millisecond * 100)
	_, err := c.Client.Logout(c.Ctx, &pb.Empty{})
	if err != nil {
		return err
	}
	c.ID = ""
	c.Player = nil
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
	if err := c.Player.JoinRoom(c.Room.ID, c.Room.OwnerSeat); err != nil {
		return err
	}
	c.Room.Players[c.Player.Seat] = c.Player
	return nil
}

func (c *MahjongClient) JoinRoom(id string) error {
	if c.Player == nil {
		return errs.ErrPlayerNotFound
	}
	if c.Player.RoomID != "" {
		return errs.ErrPlayerAlreadyInRoom
	}
	if c.Room != nil {
		return errs.ErrPlayerAlreadyInRoom
	}
	reply, err := c.Client.JoinRoom(c.Ctx, &pb.JoinRoomRequest{
		RoomID: id,
	})
	if err != nil {
		return err
	}
	roomInfo := ToRoomInfo(reply.GetRoom())
	c.Room = room.NewRoomFromInfo(roomInfo)
	if err := c.Player.JoinRoom(c.Room.ID, int(reply.Seat)); err != nil {
		return err
	}
	c.Room.Players[int(reply.Seat)] = c.Player
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
	return nil
}

func (c *MahjongClient) LeaveRoom() error {
	if c.Room == nil {
		return errs.ErrRoomNotFound
	}
	if c.ReadyStream != nil {
		if err := c.ReadyStream.Send(&pb.ReadyRequest{
			Request: &pb.ReadyRequest_LeaveRoom{
				LeaveRoom: &pb.Empty{},
			},
		}); err != nil {
			return err
		}
	}
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

func (c *MahjongClient) RefreshRoom() error {
	if c.ReadyStream == nil && c.Room == nil {
		return nil
	}
	err := c.ReadyStream.Send(&pb.ReadyRequest{
		Request: &pb.ReadyRequest_RefreshRoom{
			RefreshRoom: &pb.Empty{},
		},
	})
	return err
}

func (c *MahjongClient) StartGame() error {
	if c.Room == nil {
		return errs.ErrRoomNotFound
	}
	if c.Room.OwnerSeat != c.Player.Seat {
		return errs.ErrPlayerNotOwner
	}
	if err := c.ReadyStream.Send(&pb.ReadyRequest{
		Request: &pb.ReadyRequest_StartGame{
			StartGame: &pb.StartGameRequest{
				GameRule: ToPbGameRule(mahjong.GetDefaultRule()),
				Seed:     rand.Int63(),
			},
		},
	}); err != nil {
		return err
	}
	return nil
}

func (c *MahjongClient) RefreshGame() error {
	if c.Room == nil || c.GameStream == nil {
		return nil
	}
	return c.GameStream.Send(&pb.GameRequest{
		Request: &pb.GameRequest_RefreshGame{
			RefreshGame: &pb.Empty{},
		},
	})
}

func (c *MahjongClient) SendGameAction(action *mahjong.Call) error {
	if c.Room == nil || c.GameStream == nil {
		return nil
	}
	if action.CallType == mahjong.Next {
		c.BoardState = mahjong.NewBoardState()
	}
	return c.GameStream.Send(&pb.GameRequest{
		Request: &pb.GameRequest_Action{
			Action: ToPbCall(action),
		},
	})
}

func (c *MahjongClient) IsReady() bool {
	return c.Player != nil && c.Player.Ready
}

func (c *MahjongClient) IsOwner() bool {
	return c.Player != nil && c.Room != nil && c.Player.Seat == c.Room.OwnerSeat
}
