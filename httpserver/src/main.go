package main

import (
	// "fmt"
	"net/http"
	"time"
	"flag"
	"log"
	"context"

	// "google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/grpc"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"github.com/gin-gonic/gin"
	pb "et-protobuf3/src/gomicroapi"
	"httpserver/src/auth"
	conf "et-config/src/statusconfig"
)

type AuthResponse struct {
	Code		int 			`json:"code"`
	Msg			string 			`json:"msg"`
	Data		DataResponse 	`json:"data"`
}

type DataResponse struct {
	Token	string  `json:"token"`
}

// define token expire time
const TokenExpireDuartion = time.Hour * 2

func homehandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"code" : 2000,
		"msg"  : "success",
		"data" : gin.H{
			"username" : username,
		},
	})
}


// use gRPC call the remote Func UserLogin in tcp server
func validatePassword(userinfo pb.UserLoginInfo) (int, string) {
	// Set up a connection to the server.
	// conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials())) 	// transfer to go-micro
	etcd_reg := etcd.NewRegistry()
	// user go-micro
	service := micro.NewService(
		micro.Name("entry_task.Client"),
		micro.Registry(etcd_reg),
	)
	service.Init()

	entry_task := pb.NewTcpServerService("entry_task", service.Client())
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	// go micro 调用 UserLogin
	r, _ := entry_task.UserLogin(ctx, &pb.UserLoginInfo{Username: userinfo.Username, Password: userinfo.Password})
	log.Printf("retcode: %d\n msg:%s", r.GetRetcode(), r.GetMsg())
	return int(r.GetRetcode()), r.GetMsg()
}

// User Login API
func UserLoginHandler(c *gin.Context) {
	// 用户发送用户名，密码
	var user pb.UserLoginInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusInvalidParams,
			"msg" : conf.ErrMsg[conf.StatusInvalidParams],
		})
		return 
	}
	// 校验用户名和密码是否正确, RPC掉用
	retcode, msg := validatePassword(user)
	if retcode == 0 {
		tokenString, _ := auth.GenToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"code" : retcode,
			"msg"  : msg,
			"data" : gin.H{"token": tokenString},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : 2002,
		"msg"  : "auth failed",
	}) 
}

// User Query API
func UserQueryHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	etcd_reg := etcd.NewRegistry()
	service := micro.NewService(
		micro.Name("entry_task.Client"),
		micro.Registry(etcd_reg),
	)
	service.Init()
	entry_task := pb.NewTcpServerService("entry_task", service.Client())
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	r, err := entry_task.UserQuery(ctx, &pb.UserQueryInfo{Username: username})
	if err != nil {
		log.Println(err)
		return 
	}
	c.JSON(http.StatusOK, r)
}

func main() {
	flag.Parse()
	r := gin.Default()
	r.GET("/home", auth.JWTAuthMiddleware(), homehandler)
	r.POST("/login", UserLoginHandler)
	r.GET("/query", auth.JWTAuthMiddleware(), UserQueryHandler)
	r.Run()
}
