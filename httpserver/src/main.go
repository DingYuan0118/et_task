package main

import (
	// "fmt"
	"flag"

	// "google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/grpc"
	"httpserver/src/auth"

	hand "httpserver/src/handlerfunc"

	"github.com/gin-gonic/gin"
)

func main() {
	hand.MicrosServiceInit()
	flag.Parse()
	r := gin.Default()
	r.POST("/login", hand.UserLoginHandler)
	r.GET("/query", auth.JWTAuthMiddleware(), hand.UserQueryHandler)
	r.POST("/update-nickname", auth.JWTAuthMiddleware(), hand.UserUpdateNicknameHandler)
	r.POST("/upload-pic", auth.JWTAuthMiddleware(), hand.UserUploadPicHandler)
	r.Run()
}
