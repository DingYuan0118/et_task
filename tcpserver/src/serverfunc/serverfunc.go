package serverfunc

import (
	"context"
	"encoding/json"
	"fmt"

	// "time"

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
	rep.Retcode = conf.StatusSuccess
	rep.Msg = conf.ErrMsg[conf.StatusSuccess]

	// redis 缓存
	username := req.Username
	conn, err := rediscache.RedisInit()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		retcode := int32(conf.StatusServerError)
		msg := err.Error()
		rep.Retcode = retcode
		rep.Msg = msg
		// rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return nil
	}
	defer conn.Close()

	user := new(db.User)
	res_json, err := redis.Bytes(conn.Do("GET", username))
	// cache miss
	if err != nil {
		zaplog.Logger.Error(err.Error())
		engine, err := db.DBConnect()
		if err != nil {
			zaplog.Logger.Error(err.Error())
			retcode := int32(conf.StatusServerError)
			msg := err.Error()
			rep.Retcode = retcode
			rep.Msg = msg
			return nil
		}
		res, err := engine.Where("usr_name = ?", req.Username).Get(user)
		if err != nil {
			zaplog.Logger.Error(err.Error())
			retcode := int32(conf.StatusServerError)
			msg := err.Error()
			rep.Retcode = retcode
			rep.Msg = msg
			return nil
		}
		// 	// user not exist
		if !res {
			zaplog.Logger.Info(fmt.Sprintf("user:%s not exist", req.Username))
			rep.Retcode = conf.StatusLoginFailedNoUser
			rep.Msg = conf.ErrMsg[conf.StatusLoginFailedNoUser]
			return nil
		}

		// redis 缓存更新
		user_data_json_encode, err := json.Marshal(*user)
		if err != nil {
			zaplog.Logger.Error(err.Error())
		} else {
			// 5分钟过期
			_, err = conn.Do("set", user.Name, user_data_json_encode, "NX", "EX", "300")
			if err != nil {
				zaplog.Logger.Error(err.Error() + "redis 缓存更新时失败, username: " + username)
			}else{
				zaplog.Logger.Info("redis 缓存更新, username: " + username)
			}
		}
	} else {
		zaplog.Logger.Info("redis cache 命中, username: " + username)
		// 缓存命中，直接导入
		json.Unmarshal(res_json, user)
	}

	compareResult, _ := util.MD5ComparePasswords(user.Password, req.Password)
	if !compareResult {
		zaplog.Logger.Error("password wrong")
		retcode = conf.StatusLoginFailedPasswordWrong
		msg = conf.ErrMsg[conf.StatusLoginFailedPasswordWrong]
	} else {
		zaplog.Logger.Info("success")
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
		zaplog.Logger.Error(err.Error())
		retcode := int32(conf.StatusServerError)
		msg := err.Error()
		rep.Retcode = retcode
		rep.Msg = msg
		// rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return nil
	}
	defer conn.Close()
	user := new(db.User)
	res_json, err := redis.Bytes(conn.Do("GET", username))

	// cache miss
	if err != nil {
		zaplog.Logger.Error(err.Error())
		// read DB, return engine after sync()
		engine, err := db.DBConnect()
		if err != nil {
			zaplog.Logger.Error(err.Error())
			retcode := int32(conf.StatusServerError)
			msg := err.Error()
			rep.Retcode = retcode
			rep.Msg = msg
			return nil
		}

		res, err := engine.Where("usr_name = ?", req.Username).Get(user)
		if err != nil {
			zaplog.Logger.Error(err.Error())
			retcode := int32(conf.StatusServerError)
			msg := err.Error()
			rep.Retcode = retcode
			rep.Msg = msg
			return nil
		}
		// 理论上应该不存在 Username 不存在的场景，除非用户在查询时，用户名被删除。
		if !res {
			err := fmt.Errorf("user:%s not exist", req.Username)
			zaplog.Logger.Info(err.Error())
			rep.Retcode = conf.StatusQueryFaild
			rep.Msg = conf.ErrMsg[conf.StatusQueryFaild]
			return nil
		}
		// redis 缓存更新
		user_data_json_encode, err := json.Marshal(*user)
		if err != nil {
			zaplog.Logger.Error(err.Error())
		} else {
			// 5分钟过期
			_, err = conn.Do("set", user.Name, user_data_json_encode, "NX", "EX", "300")
			if err != nil {
				zaplog.Logger.Error(err.Error() + "redis 缓存更新时失败, username: " + username)
			}else{
				zaplog.Logger.Info("redis 缓存更新, username: " + username)
			}
		}
	} else {
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

func (s *Server) UpdateNickname(ctx context.Context, req *pb.UpdateNicknameInfo, rep *pb.UpdateNicknameReturn) error {
	// tmp := pb.UpdateNicknameReturn{}
	// return &tmp, nil
	username := req.Username
	newNickname := req.Nickname
	if len(newNickname) > 64 {
		rep.Retcode = conf.StatusNicknameTooLong
		rep.Msg = conf.ErrMsg[conf.StatusNicknameTooLong]
		return nil
	}
	
	// 先查缓存
	conn, err := rediscache.RedisInit()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		retcode := int32(conf.StatusServerError)
		msg := err.Error()
		rep.Retcode = retcode
		rep.Msg = msg
		return nil
	}
	defer conn.Close()

	// 数据库连接
	engine, err := db.DBConnect()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		retcode := int32(conf.StatusServerError)
		msg := err.Error()
		rep.Retcode = retcode
		rep.Msg = msg
		return nil
	}

	user := new(db.User)
	res_json, err := redis.Bytes(conn.Do("GET", username))
	if err != nil {
		zaplog.Logger.Error(err.Error())
		res, err := engine.Where("usr_name = ?", username).Get(user)
		if err != nil {
			zaplog.Logger.Error(err.Error())
			retcode := int32(conf.StatusServerError)
			msg := err.Error()
			rep.Retcode = retcode
			rep.Msg = msg
			return nil
		}
			// 理论上应该不存在 Username 不存在的场景，除非用户在修改时，用户名被删除。
		if !res {
			err := fmt.Errorf("user:%s not exist", req.Username)
			zaplog.Logger.Info(err.Error())
			rep.Retcode = conf.StatusUpdateNicknameFaildNoUser
			rep.Msg = conf.ErrMsg[conf.StatusUpdateNicknameFaildNoUser]
			return nil
		}
	}else{
		zaplog.Logger.Info("redis cache 命中, username: " + username)
		// 缓存命中，直接导入
		json.Unmarshal(res_json, user)
	}
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return nil
	}
	user.Nickname = newNickname

	_, err = session.Where("usr_name = ?", username).Cols("usr_nickname").Update(user)
	if err != nil {
		session.Rollback()
		zaplog.Logger.Error(err.Error())
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return nil
	}

	err = session.Commit()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return nil
	}
	zaplog.Logger.Info("mysql 数据库更新, username: " + username)

	// 缓存删除，防止不一致。 不能先更新数据再更新缓存，存在数据一致性问题。
	_, err = conn.Do("DEL", user.Name)
	if err != nil {
		zaplog.Logger.Error(err.Error())
	}else{
		zaplog.Logger.Info("redis 删除, username: " + username)
	}

	// 返回更新信息
	rep.Retcode = conf.StatusSuccess
	rep.Msg = conf.ErrMsg[conf.StatusSuccess]
	rep.Data = new(pb.UpdateNicknameReturn_Data)
	rep.Data.Nickname = user.Nickname
	return nil
}

func (s *Server) UploadPic(ctx context.Context, req *pb.UploadPicInfo, rep *pb.UploadPicReturn) error {
	// tmp := pb.UpdatePicReturn{}
	username := req.GetUsername()
	url := req.GetData().GetProfilePicUrl()
	// 先查缓存, 查缓存
	conn, err := rediscache.RedisInit()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		retcode := int32(conf.StatusServerError)
		msg := err.Error()
		rep.Retcode = retcode
		rep.Msg = msg
		// rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return nil
	}
	defer conn.Close()
	
	// 数据库连接
	engine, err := db.DBConnect()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		retcode := int32(conf.StatusServerError)
		msg := err.Error()
		rep.Retcode = retcode
		rep.Msg = msg
		return nil
	}

	user := new(db.User)

	res_json, err := redis.Bytes(conn.Do("GET", username))
	if err != nil {
		zaplog.Logger.Error(err.Error())
		res, err := engine.Where("usr_name = ?", username).Get(user)
		if err != nil {
			zaplog.Logger.Error(err.Error())
			rep.Retcode, rep.Msg = util.ThirdPackageError(err)
			return nil
		}
		// 用户不存在
		if !res {
			err := fmt.Errorf("user:%s not exist", req.Username)
			zaplog.Logger.Info(err.Error())
			rep.Retcode = conf.StatusUploadPicFailedNouser
			rep.Msg = conf.ErrMsg[conf.StatusUploadPicFailedNouser]
			return nil
		}
	}else{
		zaplog.Logger.Info("redis cache 命中, username: " + username)
		// 缓存命中，直接导入
		json.Unmarshal(res_json, user)
	}
	
	old_url := user.Profile_pic_url
	user.Profile_pic_url = url
	// 更新数据库
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return nil
	}
	_, err = engine.ID(user.Id).Cols("profile_pic_url").Update(user)
	if err != nil {
		session.Rollback()
		zaplog.Logger.Error(err.Error())
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return nil
	}
	err = session.Commit()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return nil
	}
	zaplog.Logger.Info("mysql 数据库更新, username: " + username)

	// 缓存删除，防止不一致。 不能先更新数据再更新缓存，存在数据一致性问题。
	_, err = conn.Do("DEL", user.Name)
	if err != nil {
		zaplog.Logger.Error(err.Error())
	}else{
		zaplog.Logger.Info("redis 缓存删除, username: " + username)
	}

	rep.Data = new(pb.UploadPicReturn_Data)
	rep.Data.ProfilePicUrl = url
	rep.Data.OldProfilePicUrl = old_url
	rep.Retcode = int32(conf.StatusSuccess)
	rep.Msg = conf.ErrMsg[conf.StatusSuccess]
	return nil
}
