package client

import (
	pb "github.com/hphphp123321/mahjong-client/app/api/v1"
	"github.com/hphphp123321/mahjong-client/app/errs"
	"github.com/hphphp123321/mahjong-client/app/model/player"
	"github.com/hphphp123321/mahjong-go/mahjong"
	log "github.com/sirupsen/logrus"
	"io"
	"sync"
)

func StartReadyRecvStream(c *MahjongClient, wg *sync.WaitGroup) chan error {
	done := make(chan error)
	go func() {
		for {
			reply, err := c.ReadyStream.Recv()
			if err == io.EOF {
				log.Debugln("ready stream recv EOF")
				done <- err
				return
			}
			if err != nil {
				log.Errorln("ready stream recv error:", err)
				done <- err
				break
			}
			if reply.GetMessage() != "" {
				log.Debugf("recv ready reply: %s", reply.GetMessage())
			}
			switch reply.GetReply().(type) {
			case *pb.ReadyReply_RefreshRoomReply:
				if err := handleRefreshRoomReply(c, reply, wg); err != nil {
					log.Errorln("handle refresh room reply error:", err)
				}
			case *pb.ReadyReply_GetReady:
				if err := handleGetReadyReply(c, reply, wg); err != nil {
					log.Errorln("handle get ready reply error:", err)
				}
			case *pb.ReadyReply_CancelReady:
				if err := handleCancelReadyReply(c, reply, wg); err != nil {
					log.Errorln("handle cancel ready reply error:", err)
				}
			case *pb.ReadyReply_PlayerJoin:
				if err := handlePlayerJoinReply(c, reply, wg); err != nil {
					log.Errorln("handle player join reply error:", err)
				}
			case *pb.ReadyReply_PlayerLeave:
				if err := handlePlayerLeaveReply(c, reply, wg); err != nil {
					log.Errorln("handle player leave reply error:", err)
				}
			case *pb.ReadyReply_AddRobot:
				if err := handleAddRobotReply(c, reply, wg); err != nil {
					log.Errorln("handle add robot reply error:", err)
				}
			case *pb.ReadyReply_StartGame:
				if err := handleStartGameReply(c, reply, wg); err != nil {
					log.Errorln("handle start game reply error:", err)
				}
				done <- errs.ErrGameStart
				return // start game
			case *pb.ReadyReply_Chat:
				if err := handleChatReply(c, reply, wg); err != nil {
					log.Errorln("handle chat reply error:", err)
				}
			}
		}
	}()
	return done
}

func handleRefreshRoomReply(c *MahjongClient, reply *pb.ReadyReply, wg *sync.WaitGroup) error {
	roomInfo := ToRoomInfo(reply.GetRefreshRoomReply())
	return c.Room.Refresh(roomInfo)
}

func handleGetReadyReply(c *MahjongClient, reply *pb.ReadyReply, wg *sync.WaitGroup) error {
	seat := int(reply.GetGetReady().GetSeat())
	if err := c.Room.SetReady(seat); err != nil {
		return err
	}
	if seat == c.Player.Seat && c.Player.Name == reply.GetGetReady().GetPlayerName() {
		wg.Done()
	}
	return nil
}

func handleCancelReadyReply(c *MahjongClient, reply *pb.ReadyReply, wg *sync.WaitGroup) error {
	seat := int(reply.GetCancelReady().GetSeat())
	if err := c.Room.SetUnReady(seat); err != nil {
		return err
	}
	if seat == c.Player.Seat {
		wg.Done()
	}
	return nil
}

func handlePlayerJoinReply(c *MahjongClient, reply *pb.ReadyReply, wg *sync.WaitGroup) error {
	name := reply.GetPlayerJoin().GetPlayerName()
	seat := reply.GetPlayerJoin().GetSeat()
	if name == c.Player.Name && int(seat) == c.Player.Seat {
		wg.Done()
		return nil
	}
	return c.Room.Join(player.NewPlayer(name), int(seat))
}

func handlePlayerLeaveReply(c *MahjongClient, reply *pb.ReadyReply, wg *sync.WaitGroup) error {
	seat := int(reply.GetPlayerLeave().GetSeat())
	name := reply.GetPlayerLeave().GetPlayerName()
	ownerSeat := reply.GetPlayerLeave().GetOwnerSeat()
	if name == c.Player.Name && seat == c.Player.Seat {
		if err := c.ReadyStream.CloseSend(); err != nil {
			log.Warnf("close ready stream error: %s", err)
		}
		c.Room = nil
		err := c.Player.LeaveRoom()
		if err != nil {
			log.Warnln(err)
		}
		defer func() {
			if err := recover(); err != nil {
				log.Warnf("recover from panic: %s", err)
				return
			}
		}()
		wg.Done()
		return nil
	}
	if err := c.Room.Leave(seat); err != nil {
		return err
	}
	c.Room.OwnerSeat = int(ownerSeat)
	return nil
}

func handleAddRobotReply(c *MahjongClient, reply *pb.ReadyReply, wg *sync.WaitGroup) error {
	seat := int(reply.GetAddRobot().GetRobotSeat())
	robotType := reply.GetAddRobot().GetRobotType()
	if err := c.Room.Join(player.NewRobot(robotType), seat); err != nil {
		return err
	}
	if c.IsOwner() {
		wg.Done()
	}
	return nil
}

func handleChatReply(c *MahjongClient, reply *pb.ReadyReply, wg *sync.WaitGroup) error {
	//msg := reply.GetChat().GetMessage()
	//log.Infoln(msg)
	seat := int(reply.GetChat().GetSeat())
	if seat == c.Player.Seat {
		wg.Done()
	}
	return nil
}

func handleStartGameReply(c *MahjongClient, reply *pb.ReadyReply, wg *sync.WaitGroup) error {
	if err := c.ReadyStream.CloseSend(); err != nil {
		log.Warnf("close ready stream error: %s", err)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Warnf("recover from panic: %s", err)
			return
		}
	}()
	wg.Done()
	seatsOrder := reply.GetStartGame().GetSeatsOrder()
	log.Infof("game start! order: %s, %s, %s, %s", c.Room.Players[int(seatsOrder[0])].Name, c.Room.Players[int(seatsOrder[1])].Name, c.Room.Players[int(seatsOrder[2])].Name, c.Room.Players[int(seatsOrder[3])].Name)
	c.BoardState = mahjong.NewBoardState()
	return nil
}
