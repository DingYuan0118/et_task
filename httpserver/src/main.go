package main

import (
	"errors"
	// "fmt"
	"net/http"
	"strings"
	"time"
	"flag"
	"log"
	"context"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	pb "et-protobuf3/src/rpcapi"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var hmacSampleSecret = []byte("mHpdHzQtEWQw7ntdpoNe")

// gRPC 设置
const (
	defaultName = "yuan"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

// 状态码设置
const (
	StatusSuccess = 0 					// |0|成功|
	StatusServerError = 1000			// |1000|服务器错误|
	StatusInvalidParams = 1001			// |1001|非法参数|
	StatusNotFound = 1002				// |1002|Not found|
	StatusLoginFailed = 2001			// |2001|登录失败|
	StatusTokenExpired = 2002			// |2002|Token 失效，重新登录|
	StatusQueryFaild = 3001 			// |3001|查询失败|
	StatusUpdateNicknameFaild = 3002	// |3002|更新昵称失败|
	StatusUploadPicFailed = 3003		// |3003|上传头像失败|
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type Myclaims struct {
	Username	string	`json:"username"`
	jwt.StandardClaims
}

type UserLoginInfo struct {
	Username	string		`json:"username"`
	Password	string		`json:"password"`
}

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

func GenToken(name string) (string, error) {
	claim := Myclaims{
		name, 
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuartion).Unix(),
			Issuer: "entry task",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signed_token, err := token.SignedString(hmacSampleSecret)
	return signed_token, err
}

func ParseToken(tokenString string) (*Myclaims, error){
	token, err := jwt.ParseWithClaims(tokenString, &Myclaims{}, func(token *jwt.Token) (i interface{}, err error){
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Myclaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWT 认证中间件
// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

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
func validatePasswork(userinfo UserLoginInfo) (int, string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connnect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTcpServerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UserLogin(ctx, &pb.UserLoginInfo{Username: userinfo.Username, Password: userinfo.Password})
	if err != nil {
		log.Fatalf("Func UserLogin rpc call failed: %v", err)
	}
	log.Printf("retcode: %d\n msg:%s", r.GetRetcode(), r.GetMsg())
	return int(r.GetRetcode()), r.GetMsg()
}

// User Login API
func UserLoginHandler(c *gin.Context) {
	// 用户发送用户名，密码
	var user UserLoginInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1001,
			"msg" : "invalid parameters",
		})
		return 
	}
	// 校验用户名和密码是否正确, RPC掉用
	retcode, msg := validatePasswork(user)
	if retcode == 0 {
		tokenString, _ := GenToken(user.Username)
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

func main() {
	flag.Parse()
	r := gin.Default()
	r.GET("/home", JWTAuthMiddleware(), homehandler)
	r.POST("/login", UserLoginHandler)
	r.Run()
}
