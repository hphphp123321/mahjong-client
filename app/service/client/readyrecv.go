package client

import (
	pb "github.com/hphphp123321/mahjong-client/app/api/v1"
	"github.com/hphphp123321/mahjong-client/app/model/player"
	log "github.com/sirupsen/logrus"
	"io"
)

func StartReadyRecvStream(c *MahjongClient) chan error {
	done := make(chan error)
	go func() {
		for {
			reply, err := c.ReadyStream.Recv()
			if err == io.EOF {
				done <- nil
				return
			}
			if err != nil {
				log.Errorln("ready stream recv error:", err)
				done <- err
				return
			}
			log.Debugf("recv ready reply: %s", reply.GetMessage())
			switch reply.GetReply().(type) {
			case *pb.ReadyReply_GetReady:
				if err := handleGetReadyReply(c, reply); err != nil {
					log.Errorln("handle get ready reply error:", err)
				}
			case *pb.ReadyReply_CancelReady:
				if err := handleCancelReadyReply(c, reply); err != nil {
					log.Errorln("handle cancel ready reply error:", err)
				}
			case *pb.ReadyReply_PlayerJoin:
				if err := handlePlayerJoinReply(c, reply); err != nil {
					log.Errorln("handle player join reply error:", err)
				}
			case *pb.ReadyReply_PlayerLeave:
				if err := handlePlayerLeaveReply(c, reply); err != nil {
					log.Errorln("handle player leave reply error:", err)
				}
			case *pb.ReadyReply_AddRobot:
				if err := handleAddRobotReply(c, reply); err != nil {
					log.Errorln("handle add robot reply error:", err)
				}
			case *pb.ReadyReply_StartGame:
				if err := handleStartGameReply(c, reply); err != nil {
					log.Errorln("handle start game reply error:", err)
				}
				return // start game
			case *pb.ReadyReply_Chat:
				if err := handleChatReply(c, reply); err != nil {
					log.Errorln("handle chat reply error:", err)
				}
			}
		}
	}()
	return done
}

func handleGetReadyReply(c *MahjongClient, reply *pb.ReadyReply) error {
	seat := reply.GetGetReady().GetSeat()
	return c.Room.SetReady(int(seat))
}

func handleCancelReadyReply(c *MahjongClient, reply *pb.ReadyReply) error {
	seat := reply.GetCancelReady().GetSeat()
	return c.Room.SetUnReady(int(seat))
}

func handlePlayerJoinReply(c *MahjongClient, reply *pb.ReadyReply) error {
	name := reply.GetPlayerJoin().GetPlayerName()
	seat := reply.GetPlayerJoin().GetSeat()
	if name == c.Player.Name && int(seat) == c.Player.Seat {
		return nil
	}
	return c.Room.Join(player.NewPlayer(name), int(seat))
}

func handlePlayerLeaveReply(c *MahjongClient, reply *pb.ReadyReply) error {
	seat := reply.GetPlayerLeave().GetSeat()
	name := reply.GetPlayerLeave().GetPlayerName()
	if name == c.Player.Name && int(seat) == c.Player.Seat {
		if err := c.ReadyStream.CloseSend(); err != nil {
			log.Warnf("close ready stream error: %s", err)
		}
		c.Room = nil
		return nil
	}
	if err := c.Room.Leave(int(seat)); err != nil {
		return err
	}
	return nil
}

func handleAddRobotReply(c *MahjongClient, reply *pb.ReadyReply) error {
	seat := reply.GetAddRobot().GetRobotSeat()
	robotType := reply.GetAddRobot().GetRobotType()
	return c.Room.Join(player.NewRobot(robotType), int(seat))
}

func handleChatReply(c *MahjongClient, reply *pb.ReadyReply) error {
	msg := reply.GetChat().GetMessage()
	log.Infoln(msg)
	return nil
}

func handleStartGameReply(c *MahjongClient, reply *pb.ReadyReply) error {
	if err := c.ReadyStream.CloseSend(); err != nil {
		log.Warnf("close ready stream error: %s", err)
	}
	return nil
}
