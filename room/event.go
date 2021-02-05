package room

import player "github.com/my_game/module/room"

//房间的消息类型
type Event interface {
	name()
}

//房间初始化信息
type RefreshRoomPlayersEvent struct {
	Player     *player.Player   //需要初始化的用户
	PlayerList []*player.Player //当前在用户房间的用户
}

//加入房间的消息
type JoinEvent struct {
	Player *player.Player
}

//离开房间的消息
type LeaveEvent struct {
	Player *player.Player
}

func (*RefreshRoomPlayersEvent) name()  {}
func (*JoinEvent) name()  {}
func (*LeaveEvent) name() {}
