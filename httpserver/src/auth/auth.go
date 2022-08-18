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

// Myclaims define custom messages need to be saved in token string 
type Myclaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// define token expire time
const TokenExpireDuartion = time.Minute * 5
const TotalUserNum = 10000000 // 总共一千万用户
const UserStartNum = 10000000 // 起始用户ID

// GenToken generate the token string according to the name
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

// ParseToken will validate the token, if valid return the *Myclaims struct which contained the saved message
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


// JWTAuthMiddleware is a token middleware used to verify token
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// Token 方式为 Bear Token, 置于请求头中
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
