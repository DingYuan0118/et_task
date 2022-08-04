package main

import (
	"fmt"

	// "fmt"
	// "log"
	// "net"

	"github.com/go-micro/plugins/v4/registry/etcd"
	pb "et-protobuf3/src/gomicroapi"
	"go-micro.dev/v4"
	s "tcpserver/src/serverfunc"
	// "go-micro.dev/v4/codec/proto"
	// "google.golang.org/grpc"
)

func main() {
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
	pb.RegisterTcpServerHandler(service.Server(), new(s.Server))
	err := service.Run()
	if err != nil {
		fmt.Println(err)
	}
}

