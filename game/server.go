package game

import (
	"fmt"
	"log"
	"net"
)

const (
	//最大的消息长度，默认是512字节
	MaxDataSize = 512
)

//游戏服务
func StartServer() {

	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8100,
	})

	if err != nil {
		log.Fatalln("启动game监听失败", err)
	}

	log.Println("开启game服务")
	gameRoom := NewGameRoom(socket)
	//运行房间
	go gameRoom.Run()

	//监听用户链接
	for {
		//监听连接的用户

		data := make([]byte, MaxDataSize)
		_, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("读取数据失败!", err)
			continue
		}

		//玩家加入房间
		gameRoom.JoinChan() <- remoteAddr


		//广播消息
		gameRoom.broadcast <- data
	}

}
