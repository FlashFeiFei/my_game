package main

import (
	"fmt"
	"github.com/my_game/game"
	protoc_room "github.com/my_game/module/room"
	"github.com/my_game/room"
	"google.golang.org/grpc"
	"log"
	"net"
)

//启动房间服务
func roomServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatalln("启动room监听失败", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	protoc_room.RegisterRoomGrpcServer(grpcServer, &room.RoomGrpc{})

	log.Println("启动room服务")

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("room启动失败", err)
	}

}

//启动游戏服务
func gameServer() {
	game.StartServer()
}

func main() {
	log.Println("欢迎进入")

	//启动room服务
	go roomServer()

	//启动game服务
	go gameServer()

	exit := make(chan bool)
	<-exit
}
