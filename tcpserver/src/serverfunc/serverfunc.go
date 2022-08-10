package serverfunc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	conf "et-config/src/statusconfig"
	pb "et-protobuf3/src/gomicroapi"
	"tcpserver/src/db"
	"tcpserver/src/rediscache"
	"tcpserver/src/util"
	"tcpserver/src/zaplog"

	"github.com/gomodule/redigo/redis"
)
type Server struct {
}

var retcode int32
var msg string
// var useCache bool = true

// tcp 服务端检查密码，demo
func (s *Server) UserLogin(ctx context.Context, req *pb.UserLoginInfo, rep *pb.LoginReturn) error {
	// rep.Retcode = conf.StatusSuccess
	// rep.Msg = conf.ErrMsg[conf.StatusSuccess]
	engine, err := db.DBConnect()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	user := new(db.User)
	res, err := engine.Where("usr_name = ?", req.Username).Get(user)
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	// user not exist
	if !res {
		zaplog.Logger.Info(fmt.Sprintf("user:%s not exist", req.Username))
		rep.Retcode = conf.StatusLoginFailed
		rep.Msg = conf.ErrMsg[conf.StatusLoginFailed] + "user not exist"
		return nil
	}

	compareResult, _ := util.ComparePasswords(user.Password, req.Password)
	if !compareResult {
		retcode = conf.StatusLoginFailed
		msg = conf.ErrMsg[conf.StatusLoginFailed] + "password wrong"
	}else{
		retcode = conf.StatusSuccess
		msg = conf.ErrMsg[conf.StatusSuccess]
	}

	rep.Retcode = retcode
	rep.Msg = msg
	return nil
}

func (s *Server) UserQuery(ctx context.Context, req *pb.UserQueryInfo, rep *pb.QueryReturn) error {
	// tmp := pb.QueryReturn{}
	username := req.Username
	conn, err := rediscache.RedisInit()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	defer conn.Close()
	user := new(db.User)
	res_json, err := redis.Bytes(conn.Do("GET", username))

	// cache miss
	if err != nil {
		zaplog.Logger.Info("redis cache miss, username: " + username)
		// read DB, return engine after sync()
		engine, err := db.DBConnect()
		if err != nil {
			rep.Retcode, rep.Msg = util.ThirdPackageError(err)
			return err
		}

		res, _ := engine.Where("usr_name = ?", req.Username).Get(user)
		// 理论上应该不存在 Username 不存在的场景，除非用户在查询时，用户名被删除。
		if !res {
			err := fmt.Errorf("user:%s not exist", req.Username)
			zaplog.Logger.Info(err.Error())
			rep.Retcode = conf.StatusQueryFaild
			rep.Msg = conf.ErrMsg[conf.StatusQueryFaild] + err.Error()
			return err
		}
		// redis 缓存更新
		user_data_json_encode, err := json.Marshal(*user)
		if err != nil {
			zaplog.Logger.Error(err.Error())
		}else{
			// 5分钟过期
			_, err = conn.Do("set", user.Name, user_data_json_encode, "NX", "EX", "300")
			if err != nil {
				zaplog.Logger.Error(err.Error())
			}
			zaplog.Logger.Info("redis 缓存更新, username: " + username)
		}
	}else{
		zaplog.Logger.Info("redis cache 命中, username: " + username)
		// 缓存命中，直接导入
		json.Unmarshal(res_json, user)
	}

	rep.Data = new(pb.QueryReturn_Data)
	rep.Data.Nickname = user.Nickname
	rep.Data.ProfilePic = user.Profile_pic_url
	rep.Data.Username = user.Name
	rep.Retcode = int32(conf.StatusSuccess)
	rep.Msg = conf.ErrMsg[conf.StatusSuccess]
	return nil
}

func (s *Server) UpdateNickname(ctx context.Context, req *pb.UpdateNicknameInfo, rep *pb.UpdateNicknameReturn) error{
	// tmp := pb.UpdateNicknameReturn{}
	// return &tmp, nil
	username := req.Username
	newNickname := req.Nickname
	if len(newNickname) > 64 {
		rep.Retcode = conf.StatusUpdateNicknameFaild
		rep.Msg = conf.ErrMsg[conf.StatusUpdateNicknameFaild] + "Nickname too long"
		return nil
	}
	// 更新数据库
	user := new(db.User)
	engine, err := db.DBConnect()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	res, err := engine.Where("usr_name = ?", username).Get(user)
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	// 理论上应该不存在 Username 不存在的场景，除非用户在修改时，用户名被删除。
	if !res {
		err := fmt.Errorf("user:%s not exist", req.Username)
		zaplog.Logger.Info(err.Error())
		rep.Retcode = conf.StatusUpdateNicknameFaild
		rep.Msg = conf.ErrMsg[conf.StatusUpdateNicknameFaild] + err.Error()
		return err
	}
	user.Nickname = newNickname
	// 使用事务更新 用户更新暂不考虑并发场景，不使用协程更新数据库与缓存
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	
	_, err = session.Where("usr_name = ?", username).Cols("usr_nickname").Update(user)
	if err != nil {
		session.Rollback()
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}

	err = session.Commit()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	zaplog.Logger.Info("mysql 数据库更新, username: " + username)

	// 更新缓存
	conn, err := rediscache.RedisInit()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	defer conn.Close()
	user_data_json_encode, err := json.Marshal(*user)
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	// 5分钟过期
	_, err = conn.Do("set", user.Name, user_data_json_encode, "EX", "300")
	if err != nil {
		zaplog.Logger.Error(err.Error())
	}
	zaplog.Logger.Info("redis 缓存更新, username: " + username)

	// 返回更新信息
	rep.Retcode = conf.StatusSuccess
	rep.Msg = conf.ErrMsg[conf.StatusSuccess]
	rep.Data = new(pb.UpdateNicknameReturn_Data)
	rep.Data.Nickname = user.Nickname
	return nil
}

func (s *Server) UploadPic(ctx context.Context, req *pb.UploadPicInfo, rep *pb.UploadPicReturn) error{
	// tmp := pb.UpdatePicReturn{}
	username := req.GetUsername()
	url := req.GetData().GetProfilePicUrl()
	// 更新数据库
	user := new(db.User)
	engine, err := db.DBConnect()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	
	res, err := engine.Where("usr_name = ?", username).Get(user)
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	// 用户不存在
	if !res {
		err := fmt.Errorf("user:%s not exist", req.Username)
		zaplog.Logger.Info(err.Error())
		rep.Retcode = conf.StatusUploadPicFailed
		rep.Msg = conf.ErrMsg[conf.StatusUploadPicFailed] + err.Error()
		return err
	}
	old_url := user.Profile_pic_url
	user.Profile_pic_url = url
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	_, err = engine.ID(user.Id).Cols("profile_pic_url").Update(user)
	if err != nil {
		session.Rollback()
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	err = session.Commit()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	zaplog.Logger.Info("mysql 数据库更新, username: " + username)
	// 缓存更新
	conn, err := rediscache.RedisInit()
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	defer conn.Close()
	user_data_json_encode, err := json.Marshal(*user)
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	// 5分钟过期
	_, err = conn.Do("set", user.Name, user_data_json_encode, "EX", "300")
	if err != nil {
		zaplog.Logger.Error(err.Error())
	}
	zaplog.Logger.Info("redis 缓存更新, username: " + username)
	
	rep.Data = new(pb.UploadPicReturn_Data)
	rep.Data.ProfilePicUrl = url
	rep.Data.OldProfilePicUrl = old_url
	rep.Retcode = int32(conf.StatusSuccess)
	rep.Msg = conf.ErrMsg[conf.StatusSuccess]
	return nil
}