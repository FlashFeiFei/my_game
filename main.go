package main

import (
	"fmt"
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
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	protoc_room.RegisterRoomGrpcServer(grpcServer, &room.RoomGrpc{})

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("room启动失败", err)
	}

}

func main() {
	log.Println("欢迎进入")

	//启动room服务
	go roomServer()

	exit := make(chan bool)
	<-exit
}
