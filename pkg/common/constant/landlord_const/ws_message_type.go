package landlord_const

type WsMessageType int

const (
	ReadyGame WsMessageType = iota
	UnReadyGame
	StartGame

	PlayerJoin
	PlayerExit
	Bid
	BidEnd

	PlayCard
	PleasePlayCard
	PassType

	GameEnd

	Chat

	Pong
)

func (wsMessageType WsMessageType) GetWsMessageType() string {
	return []string{"READY_GAME", "UNREADY_GAME", "START_GAME", "PLAYER_JOIN", "PLAYER_EXIT", "BID", "BID_END",
		"PLAY_CARD", "PLEASE_PLAY_CARD", "PASS", "GAME_END", "CHAT", "PONG"}[wsMessageType]
}
func (wsMessageType WsMessageType) GetWsMessageTypeName() string {
	return []string{"玩家准备", "玩家取消准备", "开始游戏", "有玩家加入", "有玩家退出", "叫牌", "叫牌结束", "有玩家出牌",
		"请出牌", "要不起", "游戏结束", "聊天消息", "心跳检测"}[wsMessageType]
}
