package main

import (
	// "fmt"
	"flag"

	// "google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/grpc"
	"httpserver/src/auth"
	"httpserver/src/zaplog"
	hand "httpserver/src/handlerfunc"

	// ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func main() {
	// 日志初始化在前
	zaplog.InitLogger()
	hand.MicrosServiceInit()
	flag.Parse()
	// r := gin.New()
	r := gin.Default()
	// 使用 zap 记录 gin 以及其他信息，默认级别为 INFO，分别输出至 ./log 与 os.stdout
	// r.Use(ginzap.Ginzap(logger, "", true), ginzap.RecoveryWithZap(logger, true))
	r.POST("/login", hand.UserLoginHandler)
	r.GET("/query", auth.JWTAuthMiddleware(), hand.UserQueryHandler)
	r.POST("/update-nickname", auth.JWTAuthMiddleware(), hand.UserUpdateNicknameHandler)
	r.POST("/upload-pic", auth.JWTAuthMiddleware(), hand.UserUploadPicHandler)

	// 支持使用 HTTP 修改日志级别
	r.GET("/log/level", gin.WrapF(zaplog.Atom.ServeHTTP))
	r.PUT("/log/level", gin.WrapF(zaplog.Atom.ServeHTTP))
	r.NoRoute(hand.NoRouteHandler)
	r.Run()
}
