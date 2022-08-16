package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	conf "et-config/src/statusconfig"
	"httpserver/src/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type Myclaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// define token expire time
const TokenExpireDuartion = time.Minute * 5
const TotalUserNum = 10000000 // 总共一千万用户
const UserStartNum = 10000000 // 起始用户ID

func GenToken(name string) (string, error) {
	claim := Myclaims{
		name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuartion).Unix(),
			Issuer:    "entry task",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signed_token, err := token.SignedString(conf.HmacSampleSecret)
	return signed_token, err
}

func ParseToken(tokenString string) (*Myclaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Myclaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return conf.HmacSampleSecret, nil
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
				"code": conf.StatusNoToken,
				"msg":  conf.ErrMsg[conf.StatusNoToken],
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		// token 格式不对
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": conf.StatusTokenInvalid,
				"msg":  conf.ErrMsg[conf.StatusTokenInvalid],
			})
			c.Abort()
			return
		}
		// token 过期
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": conf.StatusTokenInvalid,
				"msg":  conf.ErrMsg[conf.StatusTokenInvalid],
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		if mc.Username == "test" {
			num_id := util.GenerateRandomIdNum(UserStartNum, UserStartNum+TotalUserNum)
			username := "stress_test_" + strconv.Itoa(num_id)
			c.Set("username", username)
		} else {
			c.Set("username", mc.Username)
		}
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
