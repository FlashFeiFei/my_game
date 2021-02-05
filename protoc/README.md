# 前后端grpc的protobuf定义

## 功能模块

- room  游戏房间功能模块
- user 用户功能模块


### 功能模块内部命名

- bk_前缀，内部用的grpc，不对外
- rt_前缀，外部用的grpc，对外
- 其他的，自定义命名，实体命名等


## 生成grpc客户端或者服务端

- 进入功能模块,比如 cd /my_game/protoc/room
- bk_server.proto
- room.proto
- rt_server.proto

### golang DEMO
采用 protoc插件生成,其他语言采用grpc官网提供的插件生成
```cassandraql
#不分顺序,一条一条执行
#protoc   -I ./ -I ../  --go_out=./ --go-grpc_out=./ *.proto  会出问题，protoc的问题，不知道怎么解决
protoc   -I ./ -I ../  --go_out=./ --go-grpc_out=./ room.proto
protoc   -I ./ -I ../  --go_out=./ --go-grpc_out=./ room.proto
protoc   -I ./ -I ../  --go_out=./ --go-grpc_out=./ bk_server.proto
```

然后把生成的文件，拖到自己的项目里面去，能加载文件就行

命名空间一定要处理好，例如 golang的
```cassandraql
option go_package = "github.com/my_game/protoc/room";
```

### grpc生成大家都要遵循的流程

- 不同的语言中 option *_package 大家可以自己改
- 非语言差别的修改一定要商量统一

