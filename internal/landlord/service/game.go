package service

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"log"
	"math/rand"
	"time"
	"wan_go/internal/landlord/model"
	msg "wan_go/internal/landlord/model/chat_msg"
	"wan_go/internal/landlord/util"
	"wan_go/internal/landlord/ws"
	"wan_go/pkg/common/constant/landlord_const"
	"wan_go/pkg/common/db/mysql/blog"
	"wan_go/pkg/utils"
	"wan_go/sdk/service"
)

type Game struct {
	service.Service
}

func (s *Game) ReadyGame(user *blog.User) error {
	rs := Room{}
	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		return utils.Wrap(err)
	}

	player := room.GetPlayerByUserId(user.ID)
	player.Ready = true

	err = rs.UpdateRoom(room)
	if err != nil {
		return utils.Wrap(err)
	}

	nc := ws.NotifyComponent{}
	err = nc.Send2Room(room.ID, msg.NewReadyGame(user.ID))
	if err != nil {
		return utils.Wrap(err)
	}

	isAllReady := room.IsAllReady()

	if isAllReady {
		room.Mu.Lock()
		defer room.Mu.Unlock()
		s.StartGame(room)
		s.Log.Infof("StartGame %d", room.ID)
	}

	return nil
}

func (s *Game) UnReadyGame(user *blog.User) error {
	rs := Room{}
	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		return utils.Wrap(err)
	}

	player := room.GetPlayerByUserId(user.ID)
	player.Ready = false

	err = rs.UpdateRoom(room)
	if err != nil {
		return utils.Wrap(err)
	}

	nc := ws.NotifyComponent{}
	err = nc.Send2Room(room.ID, msg.NewUnReadyGame(user.ID))
	if err != nil {
		return utils.Wrap(err)
	}

	return nil
}

// Want 低于3分时，NextPlayerBid
func (s *Game) Want(user *blog.User, score int) error {
	rs := Room{}
	nc := ws.NotifyComponent{}
	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		return utils.Wrap(err)
	}

	if room.Multiple >= score {
		return errors.New("抢地主失败：比上一家叫分低")
	}
	log.Printf("[%d] 玩家 %s 叫牌，分数为 %d 分\n", room.ID, user.UserName, score)

	player := room.GetPlayerByUserId(user.ID)

	room.Multiple = score
	//记录最近一次叫分的playerId
	room.LatestBidId = player.ID

	if score > 2 {
		var landlord *blog.User
		for _, player := range room.PlayerList {
			if player.User == nil {
				return errors.New("当前玩家账号异常：player.User is nil")
			}
			if player.User.ID == user.ID {
				if player.ID != room.BiddingPlayerId {
					return errors.New("不是当前用户的叫牌回合")
				}
				landlord = player.User
				room.StepNum = player.ID
				player.Identity = landlord_const.Landlord
				player.AddCards(room.Distribution.TopCards)
			} else {
				player.Identity = landlord_const.Farmer
			}
		}
		if landlord == nil {
			return errors.New("选取的地主玩家不能为空")
		}
		room.PrePlayTime = time.Now().UnixMilli()
		err = rs.UpdateRoom(room)
		if err != nil {
			return utils.Wrap(err)
		}

		if err = nc.Send2Room(room.ID, msg.NewBidEnd()); err != nil {
			return utils.Wrap(err)
		}
		if err = nc.Send2User(landlord.ID, msg.NewPleasePlayCard()); err != nil {
			return utils.Wrap(err)
		}
		log.Printf("[%d] 玩家 %s 成为地主", room.ID, landlord.UserName)

	} else {
		nextPlayerId := player.GetNextPlayerId()
		//bid一圈
		if room.BiddingPlayerId == room.EndBidId {
			s.MustWant(room.LatestBidId, room)
			return nil
		}
		room.IncrBiddingPlayer()

		nextUser := room.GetUserByPlayerId(nextPlayerId)
		log.Printf("[%d] 玩家 %d 抢地主，分数为%d", room.ID, player.ID, score)

		if err = nc.Send2User(nextUser.ID, msg.NewBid(score)); err != nil {
			return utils.Wrap(err)
		}
	}

	return nil
}

// MustWant 叫了一圈地主 没有叫3分的情况 都没人叫地主就默认第一家是地主 叫1分
func (s *Game) MustWant(landlordPlayerId int, room *model.Room) error {
	var landlord *blog.User
	for _, player := range room.PlayerList {
		if player.ID == landlordPlayerId {
			landlord = player.User
			room.StepNum = player.ID
			player.Identity = landlord_const.Landlord
			player.AddCards(room.Distribution.TopCards)
		} else {
			player.Identity = landlord_const.Farmer
		}
	}
	if room.Multiple < 1 {
		room.Multiple = 1
	}

	room.PrePlayTime = time.Now().UnixMilli()

	rs := Room{}
	nc := ws.NotifyComponent{}

	var err error
	if err = rs.UpdateRoom(room); err != nil {
		return utils.Wrap(err)
	}
	if err = nc.Send2Room(room.ID, msg.NewBidEnd()); err != nil {
		return utils.Wrap(err)
	}
	if err = nc.Send2User(landlord.ID, msg.NewPleasePlayCard()); err != nil {
		return utils.Wrap(err)
	}
	log.Printf("[%d] 玩家 %s 成为地主", room.ID, landlord.UserName)
	return nil
}

func (s *Game) NoWant(user *blog.User) error {
	rs := Room{}
	nc := ws.NotifyComponent{}

	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		return utils.Wrap(err)
	}
	player := room.GetPlayerByUserId(user.ID)
	nextPlayerId := player.GetNextPlayerId()

	room.IncrBiddingPlayer()
	//bid一圈
	if room.BiddingPlayerId == room.EndBidId {
		landlordPlayerId := utils.IfThen(room.LatestBidId > 0, room.LatestBidId, nextPlayerId).(int)
		return s.MustWant(landlordPlayerId, room)
	}

	nextUser := room.GetUserByPlayerId(nextPlayerId)
	fmt.Printf("[%d] 玩家 %d 选择不叫，由下家 %d 玩家叫牌", room.ID, player.ID, nextPlayerId)

	return nc.Send2User(nextUser.ID, msg.NewBid(0))
}

func (s *Game) PlayCard(user *blog.User, cardList []*model.Card) (*model.RoundResult, error) {

	rs := Room{}
	nc := ws.NotifyComponent{}

	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		return nil, utils.Wrap(err)
	}

	marshal, _ := sonic.Marshal(cardList)
	log.Printf("[%d] 玩家 %s 出牌: %s", room.ID, user.UserName, string(marshal))

	player := room.GetPlayerByUserId(user.ID)

	cardType := util.GetCardsType(cardList...)
	if cardType == -1 {
		log.Printf("[%d] 玩家 %s 打出的牌不符合规则", room.ID, user.UserName)
		return nil, errors.New("玩家打出的牌不符合规则")
	}
	if room.PreCards != nil && room.PrePlayerId != player.ID {
		preType := util.GetCardsType(room.PreCards...)
		canPlay := util.CanPlayCards(cardList, room.PreCards, cardType, preType)
		if !canPlay {
			return nil, errors.New("该玩家出的牌管不了上家")
		}
	}
	removeNextPlayerRecentCards(room, player)
	player.RecentCards = cardList
	player.RemoveCards(cardList)

	playerCardMsg := msg.NewPlayCard(user, cardList, cardType)
	if err = nc.Send2Room(room.ID, playerCardMsg); err != nil {
		return nil, utils.Wrap(err)
	}

	if cardType == landlord_const.Bomb || cardType == landlord_const.JokerBomb {
		room.DoubleMultiple()
	}
	var result *model.RoundResult
	if len(player.Cards) == 0 {
		if isSpring(room, player) {
			room.DoubleMultiple()
		}
		fmt.Printf("[%d] 游戏结束，%s 获胜！", room.ID, player.Identity.GetIdentity())
		result = getResult(room, player)
		//room.Reset()
	} else {
		fmt.Printf("[%d] 玩家 %s 出牌，类型为 %s，下一个出牌者序号为：%d", room.ID,
			player.User.UserName, cardType.GetType(), player.GetNextPlayerId())
		room.PreCards = cardList
		room.PrePlayerId = player.ID
		room.IncrStep()
		nextUser := room.GetUserByPlayerId(player.GetNextPlayerId())
		if err = nc.Send2User(nextUser.ID, msg.NewPleasePlayCard()); err != nil {
			return nil, utils.Wrap(err)
		}
	}
	room.PrePlayTime = time.Now().UnixMilli()
	if err = rs.UpdateRoom(room); err != nil {
		return nil, utils.Wrap(err)
	}
	return result, nil
}

func (s *Game) PassGame(user *blog.User) error {

	rs := Room{}
	nc := ws.NotifyComponent{}

	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		return utils.Wrap(err)
	}

	player := room.GetPlayerByUserId(user.ID)
	removeNextPlayerRecentCards(room, player)
	room.IncrStep()
	room.PrePlayTime = time.Now().UnixMilli()
	if err = rs.UpdateRoom(room); err != nil {
		return utils.Wrap(err)
	}

	fmt.Printf("[%d] 玩家 %s 要不起，下一个出牌者序号为：%d", room.ID, user.UserName, player.GetNextPlayerId())
	nextUser := room.GetUserByPlayerId(player.GetNextPlayerId())
	if err = nc.Send2User(nextUser.ID, msg.NewPleasePlayCard()); err != nil {
		return utils.Wrap(err)
	}
	if err = nc.Send2Room(room.ID, msg.NewPass(user)); err != nil {
		return utils.Wrap(err)
	}
	return nil
}

func (s *Game) StartGame(room *model.Room) error {
	if room.RoomStatus == landlord_const.Playing {
		return errors.New("房间游戏已经开始")
	} else {
		room.RoomStatus = landlord_const.Playing
	}

	distribution := &model.CardDistribution{}
	room.Distribution = distribution
	distribution.Refresh()

	for _, player := range room.PlayerList {
		cards := distribution.GetCards(player.ID)
		player.AddCards(cards)
		player.Ready = true
	}

	rs := Room{}
	nc := ws.NotifyComponent{}
	var err error

	roomId := room.ID
	if err = nc.Send2Room(roomId, msg.NewStartGame(roomId)); err != nil {
		return utils.Wrap(err)
	}

	rand.Seed(time.Now().Unix())
	n := rand.Intn(3) + 1
	room.BiddingPlayerId = n
	room.EndBidId = n + 3
	player := room.GetPlayer(n)
	if err = nc.Send2User(player.User.ID, msg.NewBid(0)); err != nil {
		return utils.Wrap(err)
	}
	if err = rs.UpdateRoom(room); err != nil {
		return utils.Wrap(err)
	}

	return nil
}

// GiveCards 推荐出牌
func (s *Game) GiveCards(user *blog.User) ([][]*model.Card, error) {

	rs := Room{}

	room, err := rs.GetUserRoom(user.ID)
	if err != nil {
		return nil, utils.Wrap(err)
	}

	if room == nil || room.PreCards == nil {
		return nil, errors.New("无可推荐出牌")
	}

	player := room.GetPlayerByUserId(user.ID)

	given := util.GivePlayCards(player.Cards, room.PreCards)

	return given, nil
}

func isSpring(room *model.Room, winner *model.Player) bool {
	if winner.IsLandlord() {
		for _, player := range room.GetFarmers() {
			if len(player.Cards) < 17 {
				return false
			}
		}
		return true
	} else {
		return len(room.GetLandlord().Cards) == 17
	}
}

func getResult(room *model.Room, player *model.Player) *model.RoundResult {
	result := &model.RoundResult{
		WinIdentity: player.Identity,
		Multiple:    room.Multiple,
	}

	for _, player := range room.PlayerList {
		if player.Identity == landlord_const.Landlord {
			result.LandlordId = player.ID
		}
	}

	return result
}

func removeNextPlayerRecentCards(room *model.Room, player *model.Player) {
	nextPlayer := room.GetPlayer(player.GetNextPlayerId())
	nextPlayer.ClearRecentCards()
}
