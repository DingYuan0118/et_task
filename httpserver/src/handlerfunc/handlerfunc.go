package handlerfunc

import (
	// "fmt"
	"context"
	"net/http"
	"os"
	"time"

	conf "et-config/src/statusconfig"
	pb "et-protobuf3/src/gomicroapi"
	"httpserver/src/auth"
	"httpserver/src/util"
	"httpserver/src/zaplog"

	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/go-micro/plugins/v4/server/grpc"
	"go-micro.dev/v4"
	"go.uber.org/zap"
)

// define token expire time
var logger *zap.Logger

// user etcd register center
var service micro.Service
var entry_task pb.TcpServerService

func MicrosServiceInit() {
	// Set up a connection to the server.
	// conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials())) 	// transfer to go-micro
	etcd_reg := etcd.NewRegistry()
	// user go-micro
	service = micro.NewService(
		micro.Server(grpc.NewServer()),
		micro.Name("entry_task.Client"),
		micro.Registry(etcd_reg),
	)
	service.Init()
	entry_task = pb.NewTcpServerService("entry_task", service.Client())
	logger = zaplog.Logger
}

// use gRPC call the remote Func UserLogin in tcp server
func validatePassword(userinfo *pb.UserLoginInfo) (int, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// go micro 调用 UserLogin
	r, err := entry_task.UserLogin(ctx, &pb.UserLoginInfo{Username: userinfo.Username, Password: userinfo.Password})
	if err != nil {
		logger.Error(err.Error())
		return conf.StatusServerError, err.Error()
	}
	if r.Retcode != 0 {
		logger.Error(r.GetMsg())
	}
	return int(r.GetRetcode()), r.GetMsg()
}

// User Login API
func UserLoginHandler(c *gin.Context) {
	// 用户发送用户名，密码
	var user pb.UserLoginInfo
	err := c.ShouldBind(&user)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusInvalidParams,
			"msg":  conf.ErrMsg[conf.StatusInvalidParams],
		})
		return
	}
	// 校验用户名和密码是否正确, RPC调用
	retcode, msg := validatePassword(&user)
	if retcode == 0 {
		tokenString, _ := auth.GenToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"code": retcode,
			"msg":  msg,
			"data": gin.H{"token": tokenString},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": retcode,
		"msg":  msg,
	})
}

// User Query API
func UserQueryHandler(c *gin.Context) {
	// 由 token 提取 username
	username := c.MustGet("username").(string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10) // 请求超时时长一般设置 10s
	defer cancel()

	r, err := entry_task.UserQuery(ctx, &pb.UserQueryInfo{Username: username})
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusServerError,
			"msg":  err.Error(),
		})
		return
	}
	if r.GetRetcode() != 0 {
		logger.Error(r.GetMsg())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": r.GetRetcode(),
		"msg":  r.GetMsg(),
		"data": gin.H{
			"username":    r.GetData().GetUsername(),
			"nickname":    r.GetData().GetNickname(),
			"profile_pic": r.GetData().GetProfilePic(),
		},
	})
}

// update Nickname API
func UserUpdateNicknameHandler(c *gin.Context) {
	var user pb.UpdateNicknameInfo
	// 由 JWT 中间件设置 username
	username := c.MustGet("username").(string)
	err := c.ShouldBind(&user)
	// 使用 Token 中的 username 防止误修改其他用户信息
	user.Username = username
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusInvalidParams,
			"msg":  conf.ErrMsg[conf.StatusInvalidParams],
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r, err := entry_task.UpdateNickname(ctx, &user)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusServerError,
			"msg":  err.Error(),
		})
		return
	}
	if r.GetRetcode() != 0 {
		logger.Error(r.GetMsg())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": r.GetRetcode(),
		"msg":  r.GetMsg(),
		"data": gin.H{
			"nickname": r.GetData().GetNickname(),
		},
	})
}

// Upload pic api
func UserUploadPicHandler(c *gin.Context) {
	var user pb.UploadPicInfo
	// username := c.PostForm("username")
	// 从 Token 中获取
	username := c.MustGet("username").(string)
	user.Username = username
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error(err.Error())
	}
	file_os, _ := file.Open()
	filetype, err := util.GetFileContentType(file_os)
	logger.Info("upload file type",
		zap.String("file type", filetype))
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusServerError,
			"msg":  err.Error(),
		})
		return
	}

	// 判断是否符合格式
	var types = []string{"image/png", "image/jpeg", "image/bmp"}
	if !util.Contains(types, filetype) {
		logger.Error("file format error, support [png, jpeg, bmp], got " + filetype)
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusUploadPicFormatWrong,
			"msg":  conf.ErrMsg[conf.StatusUploadPicFormatWrong],
		})
		return
	}
	// filesize should < 3MB
	if file.Size > 3<<20 {
		logger.Error("file too large, should less than 3MB")
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusUploadPicTooLarge,
			"msg":  conf.ErrMsg[conf.StatusUploadPicTooLarge],
		})
		return
	}
	_, err = os.Stat(conf.ImageFolder)
	if !os.IsExist(err) {
		os.MkdirAll(conf.ImageFolder, 0777) // 创建文件夹
	}
	url := conf.ImageFolder + user.Username + "." + filetype[6:]
	user.Data = new(pb.UploadPicInfo_Data)
	user.Data.ProfilePicUrl = url

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 返回旧 url 用于删除
	r, err := entry_task.UploadPic(ctx, &user)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusServerError,
			"msg":  err.Error(),
		})
		return
	}

	if r.GetRetcode() != 0 {
		logger.Error(r.GetMsg())
		c.JSON(http.StatusOK, gin.H{
			"code": r.GetRetcode(),
			"msg":  r.GetMsg(),
		})
		return
	}

	old_url := r.Data.OldProfilePicUrl
	// 删除后更新
	err = os.Remove(old_url)
	if err != nil {
		logger.Error(err.Error())
	}

	// 可能出现权限问题
	err = c.SaveUploadedFile(file, url)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": conf.StatusServerError,
			"msg":  err.Error(),
		})
		return
	}
	logger.Info("upload pic success")
	c.JSON(http.StatusOK, gin.H{
		"code": r.GetRetcode(),
		"msg":  "upload pic success",
		"data": gin.H{
			"profile_pic_url": r.GetData().GetProfilePicUrl(),
		},
	})
}

func NoRouteHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": conf.StatusNotFound,
		"msg":  conf.ErrMsg[conf.StatusNotFound],
	})
}
