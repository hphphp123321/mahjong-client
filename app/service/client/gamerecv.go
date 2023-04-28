package client

import (
	pb "github.com/hphphp123321/mahjong-client/app/api/v1"
	"github.com/hphphp123321/mahjong-client/app/errs"
	"github.com/hphphp123321/mahjong-go/mahjong"
	log "github.com/sirupsen/logrus"
	"io"
)

func StartGameRecvStream(c *MahjongClient, actionChan chan mahjong.Calls) chan error {
	done := make(chan error)
	go func() {
		for {
			reply, err := c.GameStream.Recv()
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
			if reply.GetReply() == nil {
				switch reply.GetReply().(type) {
				case *pb.GameReply_RefreshGameReply:
					if err := handleRefreshGameReply(c, reply); err != nil {
						log.Errorln("handle refresh game reply error:", err)
					}
				}
			}
			if reply.GetEvents() != nil {
				log.Debugln("recv game events:", reply.GetEvents())
				if err := handleGameEvents(c, reply); err != nil {
					log.Errorln("handle game events error:", err)
				}
			}
			if reply.GetValidActions() != nil {
				log.Debugln("recv valid actions:", reply.GetValidActions())
				validActions := ToMahjongCalls(reply.GetValidActions())
				actionChan <- validActions
			}
			if reply.GetEnd() {
				log.Debugln("recv game end")
				done <- errs.ErrGameEnd
				return
			}
		}
	}()
	return done
}

func handleRefreshGameReply(c *MahjongClient, reply *pb.GameReply) error {
	b := ToMahjongBoardState(reply.GetRefreshGameReply())
	if !b.Equal(c.BoardState) {
		log.Warnln("board state not equal")
	}
	return nil
}

func handleGameEvents(c *MahjongClient, reply *pb.GameReply) error {
	events := reply.GetEvents()
	if len(events) == 0 {
		return errs.ErrEventsEmpty
	}
	es := ToMahjongEvents(events)
	c.BoardState.DecodeEvents(es)
	return nil
}
