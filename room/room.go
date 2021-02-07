package room

import (
	protoc_room "github.com/my_game/module/room"
	"log"
)

func init() {
	if room == nil {
		room = NewRoom()

		//房间跑起来
		go room.Run()
	}
}

var room *Room

//房间
type Room struct {
	MaxCount  uint32           //房间最大人数
	Status    uint32           ///状态,0:空闲中，1:满人,2:游戏中,不怎么会出现资源竞争，不需要加锁
	playerMap map[*Client]bool //加入房间的client
	join      chan *Client     //加入房间,  无缓冲，一个一个进入
	leave     chan *Client     //退出房间,有缓冲
	broadcast chan Event       //广播消息,无缓冲，要保证事件的顺序
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

//空闲和满人之间切换
func (r *Room) updataStatus() {
	if r.Status == 2 {
		//游戏中的状态不能切换空闲和满人
		return
	}

	if r.IsMaxCount() {
		//满人
		r.Status = 1
	} else {
		//空闲中
		r.Status = 0
	}
}

//是否满人
func (r *Room) IsMaxCount() bool {

	if len(r.playerMap) < int(r.MaxCount) {
		return false
	}

	return true
}

//运行房间
func (r *Room) Run() {
	for {

		r.updataStatus() //更新房间是空闲还是满人状态
		log.Println("此时房间存活的人", r.playerMap)

		select {

		case player := <-r.join:
			//有用户加入房间
			r.playerMap[player] = true

		case leave := <-r.leave:
			//有用户离开房间
			if _, ok := r.playerMap[leave]; ok {
				//房间取消这个用户
				delete(r.playerMap, leave)
				//关闭这个用户相关的资源
				leave.Close()

				go func() {
					//广播一个用户离开的事件，通知其他用户刷新房间存活的用户，用异步，不然会死锁
					//到底是用异步好呢，还是用有缓冲的chan好呢
					r.broadcast <- &LeaveEvent{Player: &protoc_room.Player{
						User: &protoc_room.User{
							Id: leave.GetPlayer().User.Id,
						},
					}}
				}()
			}


		case event := <-r.broadcast:
			//广播资源,广播

			switch event.(type) {

			case *LeaveEvent:
				//离开事件
				//服务端监听到某个链接异常断开的时候，给存活的其他链接推送一个离开事件，让客户端调用刷新RefreshRoomPlayersEvent事件
				//来刷新房间中的用户
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
