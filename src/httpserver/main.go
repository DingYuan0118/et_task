package main

import (
	"errors"
	// "fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var hmacSampleSecret = []byte("mHpdHzQtEWQw7ntdpoNe")


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

type UserInfo struct {
	Username	string
	Password	string
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

func authHandler(c *gin.Context) {
	// 用户发送用户名，密码
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1001,
			"msg" : "无效的参数",
		})
		return 
	}
	//TODO
	// 校验用户名和密码是否正确, RPC掉用
	if true {
		tokenString, _ := GenToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"code" : 0,
			"msg"  : "success",
			"data" : gin.H{"token": tokenString},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code" : 2002,
		"msg"  : "鉴权失败",
	}) 
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

func main() {
	r := gin.Default()
	r.GET("/home", JWTAuthMiddleware(), homehandler)
	r.POST("/auth", authHandler)
	r.Run()
}
