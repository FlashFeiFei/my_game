package main

import (
	"context"
	protoc_room "github.com/my_game/module/room"
	"google.golang.org/grpc"
	"log"
	"time"
)

//获取链接
func getConn() *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:8080", opts...)
	if err != nil {
		log.Fatalln("拨号失败", err)
	}

	return conn
}

func main() {

	conn := getConn()
	client := protoc_room.NewRoomGrpcClient(conn)

	roomStreamClient, err := client.Room(context.Background())

	if err != nil {
		log.Fatalln("连接失败", err)
	}

	//用户
	userId := time.Now().Unix()
	user := &protoc_room.User{
		Id: uint64(userId),
	}

	log.Println("我的数据",user)

	//写入操作

	//发送加入房间事件
	joinEvent := &protoc_room.RoomStreamRequest_JoinEvent{
		JoinEvent: &protoc_room.JoinEvent{
			Room: nil,
			Player: &protoc_room.Player{
				User: user,
			},
		},
	}

	joinData := &protoc_room.RoomStreamRequest{Event: joinEvent}

	err = roomStreamClient.Send(joinData)
	if err != nil {
		log.Fatalln("加入房间失败", err)
	}

	log.Println("发送加入房间事件成功")

	//发送获取房间事件
	refreshRoomPlayersEvent := &protoc_room.RoomStreamRequest_RoomPlayersEvent{
		RoomPlayersEvent: &protoc_room.RefreshRoomPlayersEvent{
			Room:       nil,
			PlayerList: nil,
			Player: &protoc_room.Player{
				User: user,
			},
		},
	}

	refreshData := &protoc_room.RoomStreamRequest{Event: refreshRoomPlayersEvent}

	err = roomStreamClient.Send(refreshData)
	if err != nil {
		log.Fatalln("获取数据失败", err)
	}

	log.Println("发送获取房间事件成功")

	//读取数据
	go func() {

		for {

			in, err := roomStreamClient.Recv()

			if err != nil {
				log.Fatalln("读入失败", err)
			}

			//打印数据
			switch event := in.Event.(type) {
			case *protoc_room.RoomStreamResponse_RoomPlayersEvent:
				log.Println("房间中的用户", event.RoomPlayersEvent)
			case *protoc_room.RoomStreamResponse_LeaveEvent:
				log.Println("有用户离开房间了", event.LeaveEvent)
			}
		}

	}()

	log.Println("ctr+c 退出程序")
	stop := make(chan bool, 1)
	<-stop

	log.Println("退成成功，你能不能行，还没写好？")
}
