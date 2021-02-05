package game

import (
	"fmt"
	"net"
)

//玩家
type Player struct {
	Id string //玩家id
}

//游戏房间
var Game *game

func init() {
	if Game == nil {
		Game = &game{
			status:        0,
			loadPlayerMap: nil,
			playerConnMap: nil,
		}
	}
}

//游戏
type game struct {
	status        int                     //0空闲中，1游戏启动加载用户中，2游戏开始进行中
	loadPlayerMap map[string]*Player      //加载参与游戏的玩家
	playerConnMap map[string]*net.UDPAddr //存活玩家的udp连接
}

//设置加载初始化的用户
func (g *game) SetLoadPlayerList(loadPlayerMap map[string]*Player) {
	g.loadPlayerMap = loadPlayerMap
}

//结束游戏
//重置属性
func (g *game) EndGame() {
	g.status = 0
	g.loadPlayerMap = nil
	g.playerConnMap = nil
}

//设置玩家的udp链接
//player玩家
//conn 链接
func (g *game) SetPlayerUdpConn(player *Player, conn *net.UDPAddr) {
	if g.playerConnMap == nil {
		g.playerConnMap = make(map[string]*net.UDPAddr)
	}

	g.playerConnMap[player.Id] = conn
}

//删除玩家的udp链接
func (g *game) DeletePlayerUdpConn(player *Player) {
	if g.playerConnMap == nil {
		return
	}
	delete(g.playerConnMap, player.Id)
	return
}

//启动游戏
//游戏从空闲中，切换到游戏加载中
func (g *game) StartGame() error {

	//游戏不是空闲中不能加载
	if g.status != 0 {
		var errMsg string
		switch g.status {
		case 1:
			errMsg = "游戏启动加载中"
		case 2:
			errMsg = "游戏开始中"
		}

		return fmt.Errorf(errMsg)
	}

	if g.loadPlayerMap == nil || len(g.loadPlayerMap) <= 0 {
		return fmt.Errorf("加载用户为空")
	}

	//切换游戏到加载中
	g.status = 1

	return nil
}
