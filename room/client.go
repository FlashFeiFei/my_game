package room

import (
	protoc_room "github.com/my_game/module/room"
	"io"
)

//玩家
type Client struct {
	room   *Room                           //玩家所在的房间
	send   chan Event                      //通知给玩家的消息,有缓冲通道
	conn   protoc_room.RoomGrpc_RoomServer //玩家的链接
	player protoc_room.Player              //玩家信息
}

//创建一个玩家
func NewClient(room *Room, conn protoc_room.RoomGrpc_RoomServer) *Client {
	return &Client{
		room: room,
		send: make(chan Event, 20),
		conn: conn,
	}
}

//发送消息给客户
func (c *Client) Send() chan<- Event {
	return c.send
}

//设置玩家的基本信息
func (c *Client) SetPlayer(player protoc_room.Player) {
	c.player = player
}

//获取玩家的基本信息
func (c *Client) GetPlayer() protoc_room.Player {
	return c.player
}

//关闭玩家的资源
func (c *Client) Close() {
	close(c.send) //关闭玩家开启的send chan
}

//读取客户数据，写入房间
func (c *Client) ReadPump() {

	//链接异常时候，处理一下房间内的用户
	defer func() {
		//释放用户
		c.room.Leave() <- c
		//关闭客户资源
		c.Close()
	}()

	for {
		inputMsg, err := c.conn.Recv()

		if err == io.EOF {
			//读取完毕了，没有数据了,继续监听
			continue
		}

		if err != nil {
			//链接出问题，直接，结束读取方法就好
			return
		}

		switch event := inputMsg.Event.(type) {

		case *protoc_room.RoomStreamRequest_JoinEvent:
			//加入事件,只通知房间，不进行任何广播
			c.room.Join() <- c //玩家加入房间

		case *protoc_room.RoomStreamRequest_RoomPlayersEvent:
			//初始化事件
			//通知房间，xxx用户需要获取初始化事件
			c.room.AddBroadcast() <- &RefreshRoomPlayersEvent{
				Player:     event.RoomPlayersEvent.Player,
				PlayerList: nil,
			}

		case *protoc_room.RoomStreamRequest_LeaveEvent:
			//离开事件,只通知房间，不进行任何的广播
			c.room.Leave() <- c
		}

	}
}

//读取房间数据，发送给客户
func (c *Client) WriteRump() {
	defer func() {
		//释放用户
		c.room.Leave() <- c

		//关闭客户资源
		c.Close()
	}()

	for {
		select {
		case event, ok := <-c.send:
			if !ok {
				//有缓冲的chan，需要判断是否读取ok不，不ok说明chan被关闭了
				//函数结束
				return
			}

			err := c.conn.Send(c.buildRoomStreamResponse(event))

			if err != nil {
				//连接出了问题,结束方法
				return
			}

			//将客户的缓冲消息也同步过去
			n := len(c.send)
			for i := 0; i < n; i++ {
				err := c.conn.Send(c.buildRoomStreamResponse(<-c.send))
				if err != nil {
					//链接出了问题，结束用户
					return
				}
			}
		}
	}
}

func (c *Client) buildRoomStreamResponse(event Event) *protoc_room.RoomStreamResponse {

	roomStreamResponse := new(protoc_room.RoomStreamResponse)

	switch event.(type) {
	case *RefreshRoomPlayersEvent:
		//刷新房间用户事件
		refreshRoomPlayersEvent := event.(*RefreshRoomPlayersEvent)
		roomStreamResponse.Event = &protoc_room.RoomStreamResponse_RoomPlayersEvent{
			RoomPlayersEvent: &protoc_room.RefreshRoomPlayersEvent{
				Room:       nil,
				PlayerList: refreshRoomPlayersEvent.PlayerList,
				Player:     refreshRoomPlayersEvent.Player,
			},
		}
	}

	return roomStreamResponse
}
