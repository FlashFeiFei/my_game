package room

import (
	protoc_room "github.com/my_game/module/room"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type RoomGrpc struct {
	protoc_room.UnimplementedRoomGrpcServer
}

//room房间的grpc
func (grpc *RoomGrpc) Room(conn protoc_room.RoomGrpc_RoomServer) error {
	if room.IsMaxCount() {
		return status.Errorf(codes.Unknown, "房间已经满人了，请等待")
	}

	client := NewClient(room, conn)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {

		defer wg.Done()

		//读取用户数据，写入房间
		client.ReadPump()
	}()

	go func() {
		defer wg.Done()

		//读取房间数据，写入客户
		client.WriteRump()
	}()

	wg.Wait()

	return nil
}
