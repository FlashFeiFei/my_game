package game

import (
	"log"
	"net"
)

type GameRoom struct {
	conn       *net.UDPConn            //udp服务器的唯一链接
	udpAddrMap map[string]*net.UDPAddr //加入房间的client
	join       chan *net.UDPAddr       //加入房间,  无缓冲，一个一个进入
	broadcast  chan []byte             //广播消息,无缓冲，要保证事件的顺序
}

//创建一个房间
func NewGameRoom(conn *net.UDPConn) *GameRoom {
	return &GameRoom{
		conn:       conn,
		udpAddrMap: make(map[string]*net.UDPAddr),
		join:       make(chan *net.UDPAddr),
		broadcast:  make(chan []byte),
	}
}

//加入房间
func (g *GameRoom) JoinChan() chan<- *net.UDPAddr {
	return g.join
}


//运行房间
func (g *GameRoom) Run() {
	for {
		select {
		case udpAddr := <-g.join:
			//加入房间
			g.udpAddrMap[udpAddr.String()] = udpAddr

		case data := <-g.broadcast:

			log.Println("当前存活连接",g.udpAddrMap)
			//广播消息
			for _, udpAddr := range g.udpAddrMap {
				n, err := g.conn.WriteToUDP(data, udpAddr)
				if err != nil || n == 0{
					//错误或者,写了0行数据?
					//udp是无连接的，监听不到连接是否断开，所以当
					continue
				}
			}
		}
	}
}
