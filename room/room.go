package room

import protoc_room "github.com/my_game/module/room"

//房间
type Room struct {
	MaxCount  uint32           //房间最大人数
	Status    uint32           ///状态,0:空闲中，1:满人,2:游戏中,不怎么会出现资源竞争，不需要加锁
	playerMap map[*Client]bool //加入房间的client
	join      chan *Client     //加入房间,  无缓冲，一个一个进入
	leave     chan *Client     //退出房间,有缓冲
	broadcast chan Event       //广播消息
}

//创建房间
func NewRoom() *Room {
	return &Room{
		MaxCount:  10,
		Status:    0, //空闲中
		playerMap: make(map[*Client]bool),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan Event),
	}
}

//加入房间
func (r *Room) Join() chan<- *Client {
	return r.join
}

//离开房间
func (r *Room) Leave() chan<- *Client {
	return r.leave
}

//消息广播
func (r *Room) AddBroadcast() chan<- Event {
	return r.broadcast
}

//房间当前人数
func (r *Room) CurrentCount() int {
	return len(r.playerMap)
}

//运行房间
func (r *Room) Run() {
	for {
		select {

		case player := <-r.join:
			//有用户加入房间
			r.playerMap[player] = true

		case leave := <-r.leave:
			//有用户离开房间
			if _, ok := r.playerMap[leave]; ok {
				//房间取消这个用户
				delete(r.playerMap, leave)
				//关闭的资源
				leave.Close()
			}

		case event := <-r.broadcast:
			//广播资源,广播

			switch event.(type) {
			case *RefreshRoomPlayersEvent:
				//只通知需要初始化房间的玩家
				//构建房间内所有的玩家数据
				playerList := make([]*protoc_room.Player, 0)
				for roomPlayer, _ := range r.playerMap {
					player := roomPlayer.GetPlayer()
					playerList = append(playerList, &player)
				}

				event.(*RefreshRoomPlayersEvent).PlayerList = playerList

			}

			//信息广播
			for player := range r.playerMap {
				select {
				case player.Send() <- event: //发送消息给客户
				default:
					//没有任何一个chan匹配，发送失败，说明客户的链接异常
					//房间取消这个用户
					delete(r.playerMap, player)
					//关闭的资源
					player.Close()
				}
			}
		}
	}
}