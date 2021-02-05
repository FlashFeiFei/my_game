package room

import (
	protoc_room "github.com/my_game/module/room"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoomGrpc struct {
	protoc_room.UnimplementedRoomGrpcServer
}

//room房间的grpc
func (grpc *RoomGrpc) Room(client protoc_room.RoomGrpc_RoomServer) error {
	return status.Errorf(codes.Unimplemented, "method Room not implemented")
}
