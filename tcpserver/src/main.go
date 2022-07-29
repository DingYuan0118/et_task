package main

import (
	"context"
	"fmt"

	// "fmt"
	// "log"
	// "net"

	"github.com/go-micro/plugins/v4/registry/etcd"
	pb "et-protobuf3/src/gomicroapi"
	"go-micro.dev/v4"
	// "go-micro.dev/v4/codec/proto"
	// "google.golang.org/grpc"
)

// var (
// 	port = flag.Int("port", 50051,  "The server port")
// )
type server struct {
	// pb.UnimplementedTcpServerServer // 转为使用 go-micro
}

// tcp 服务端检查密码，demo
func (s *server) UserLogin(ctx context.Context, req *pb.UserLoginInfo, rep *pb.LoginReturn) error {
	rep.Msg = "login success"
	rep.Retcode = 0
	return nil
}

func (s *server) UserQuery(ctx context.Context, req *pb.UserQueryInfo, rep *pb.QueryReturn) error {
	// tmp := pb.QueryReturn{}
	return nil
}

func (s *server) UpdateNickname(ctx context.Context, in *pb.UpdateNicknameInfo, rep *pb.UpdateNicknameReturn) error{
	// tmp := pb.UpdateNicknameReturn{}
	// return &tmp, nil
	return nil
}

func (s *server) UploadPic(ctx context.Context, in *pb.UploadPicInfo, rep *pb.UpdatePicReturn) error{
	// tmp := pb.UpdatePicReturn{}
	return nil
}

func main() {
	// flag.Parse()
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// s := grpc.NewServer()
	// pb.RegisterTcpServerServer(s, &server{})
	// log.Printf("server listening at %v", lis.Addr())
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }

	// use etcd 作为注册存储中心
	etcd_reg := etcd.NewRegistry()
	// user go-micro
	// 使用 micro 框架实现服务注册
	service := micro.NewService(
		micro.Name("entry_task"),
		micro.Registry(etcd_reg),
	)

	// 初始化，解析命令行参数
	service.Init()
	
	// 注册服务
	pb.RegisterTcpServerHandler(service.Server(), new(server))
	
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

