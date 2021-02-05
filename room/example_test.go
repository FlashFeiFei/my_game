package room

import (
	"context"
	"fmt"
	protoc_room "github.com/game/room/protoc/room"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"log"
	"testing"
)

//加入房间测试
func TestJoinRoom(t *testing.T) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:8080", opts...)
	if err != nil {
		log.Fatalln("链接失败", err)
	}

	defer conn.Close()

	client := protoc_room.NewRoomServiceClient(conn)
	for i := 0; i < 10; i++ {
		response, err := client.Join(context.Background(), &protoc_room.JoinRequest{
			Player: &protoc_room.Player{
				Id: fmt.Sprintf("ID号:%d", i),
			},
		})

		log.Println("响应结果", response, err)
	}

}

//离开房间测试
func TestLeaveRoom(t *testing.T) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:8080", opts...)
	if err != nil {
		log.Fatalln("链接失败", err)
	}

	defer conn.Close()

	client := protoc_room.NewRoomServiceClient(conn)
	response, err := client.Leave(context.Background(), &protoc_room.LeaveRequest{
		Player: &protoc_room.Player{
			Id: fmt.Sprintf("ID号:%d", 6),
		},
	})

	log.Println("响应结果", response, err)
}

//获取房间信息测试
func TestBeginGame(t *testing.T) {

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:8080", opts...)
	if err != nil {
		log.Fatalln("链接失败", err)
	}

	defer conn.Close()

	client := protoc_room.NewRoomServiceClient(conn)
	response, err := client.BeginGame(context.Background(), &emptypb.Empty{})

	log.Println("响应结果", response, err)
}

//结束游戏
func TestEndGame(t *testing.T) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial("127.0.0.1:8080", opts...)
	if err != nil {
		log.Fatalln("链接失败", err)
	}

	defer conn.Close()

	client := protoc_room.NewRoomServiceClient(conn)
	response, err := client.EndGame(context.Background(), &emptypb.Empty{})

	log.Println("响应结果", response, err)
}
