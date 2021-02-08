package game

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8100,
	})
	if err != nil {
		fmt.Println("连接失败!", err)
		return
	}

	id := time.Now().Unix()

	//发送一个数据
	// 发送数据
	senddata := []byte(fmt.Sprintf("发送第一个数据%d", id))
	_, err = socket.Write(senddata)
	if err != nil {
		log.Fatalln("发送第一个数据失败")
	}
	log.Println("发送成功")

	go func() {
		for {
			// 发送数据
			senddata := []byte(fmt.Sprintf("定时发送数据%d", id))
			_, err = socket.Write(senddata)
			if err != nil {
				log.Fatalln("写数据失败链接异常")
			}

		}
	}()

	// 接收数据
	for {
		data := make([]byte, MaxDataSize)
		_, _, err := socket.ReadFromUDP(data)
		if err != nil {
			log.Fatalln("读取数据失败链接异常", err)

		}
		fmt.Println(string(data))
	}
}
