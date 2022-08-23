package main

import (
	"fmt"

	s "tcpserver/src/serverfunc"

	pb "et-protobuf3/src/gomicroapi"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
)

func main() {
	// use etcd 作为注册存储中心
	etcd_reg := etcd.NewRegistry()
	// user go-micro
	// 使用 micro 框架实现服务注册
	service := micro.NewService(
		micro.Server(grpc.NewServer()),
		micro.Name("entry_task"),
		micro.Registry(etcd_reg),
	)
	//
	// zaplog.InitLogger() // 在包有初始化函数场景下，该语句多余
	// 初始化，解析命令行参数
	service.Init()
	// 注册服务
	pb.RegisterTcpServerHandler(service.Server(), new(s.Server))
	err := service.Run()
	if err != nil {
		fmt.Println(err)
	}
}
