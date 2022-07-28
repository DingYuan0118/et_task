package main

import (
	"flag"
	"context"
	"net"
	"fmt"
	"log"

	"google.golang.org/grpc"
	pb "entry_task/src/rpcapi"
)

var (
	port = flag.Int("port", 50051,  "The server port")
)
type server struct {
	pb.UnimplementedTcpServerServer
}

// tcp 服务端检查密码，demo
func (s *server) UserLogin(ctx context.Context, in *pb.UserLoginInfo) (*pb.LoginReturn, error) {
	var loginreturn = pb.LoginReturn{
		Retcode: 0,
		Msg: "auth success",
	}
	return &loginreturn, nil
}

func (s *server) UserQuery(ctx context.Context, in *pb.UserQueryInfo) (*pb.QueryReturn, error){
	tmp := pb.QueryReturn{}
	return &tmp, nil
}

func (s *server) UpdateNickname(ctx context.Context, in *pb.UpdateNicknameInfo) (*pb.UpdateNicknameReturn, error){
	tmp := pb.UpdateNicknameReturn{}
	return &tmp, nil
}

func (s *server) UploadPic(ctx context.Context, in *pb.UploadPicInfo) (*pb.UpdatePicReturn, error){
	tmp := pb.UpdatePicReturn{}
	return &tmp, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTcpServerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}