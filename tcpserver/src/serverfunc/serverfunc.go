package serverfunc

import (
	"context"
	"encoding/json"
	"log"

	conf "et-config/src/statusconfig"
	pb "et-protobuf3/src/gomicroapi"
	"tcpserver/src/db"
	"tcpserver/src/rediscache"
	"tcpserver/src/util"

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
	engine, err := db.DBConnect(db.DBname, db.Password)
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	defer engine.Close()
	user := new(db.User)
	res, err := engine.Where("usr_name = ?", req.Username).Get(user)
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	if !res {
		rep.Retcode = conf.StatusLoginFailed
		rep.Msg = conf.ErrMsg[conf.StatusLoginFailed]
		return nil
	}

	compareResult, _ := util.ComparePasswords(user.Password, req.Password)
	if !compareResult {
		retcode = conf.StatusLoginFailed
		msg = conf.ErrMsg[conf.StatusLoginFailed]
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
		log.Println("redis cache miss")
		// read DB, return engine after sync()
		engine, err := db.DBConnect(db.DBname, db.Password)
		if err != nil {
			rep.Retcode, rep.Msg = util.ThirdPackageError(err)
			return err
		}
		defer engine.Close()

		res, err := engine.Where("usr_name = ?", req.Username).Get(user)
		if err != nil {
			rep.Retcode, rep.Msg = util.ThirdPackageError(err)
			return err
		}
		// 理论上应该不存在 Username 不存在的场景，除非用户在查询时，用户名被删除。 为测试补充
		if !res {
			rep.Retcode = conf.StatusQueryFaild
			rep.Msg = conf.ErrMsg[conf.StatusQueryFaild]
			return nil
		}

		// 数据 encode 准备打入 redis
		user_data_json_encode, err := json.Marshal(*user)
		if err != nil {
			log.Println(err)
		}else{
			// 5分钟过期
			_, err = conn.Do("set", user.Name, user_data_json_encode, "NX", "EX", "300")
			if err != nil {
				log.Println(err)
			}
		}
	}else{
		log.Println("redis cache 命中")
		// 缓存命中，直接导入
		err = json.Unmarshal(res_json, user)
		if err != nil {
			rep.Retcode, rep.Msg = util.ThirdPackageError(err)
			return err
		}
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
		rep.Msg = conf.ErrMsg[conf.StatusUpdateNicknameFaild]
		return nil
	}
	// 更新数据库
	user := new(db.User)
	engine, err := db.DBConnect(db.DBname, db.Password)
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	defer engine.Close()
	res, err := engine.Where("usr_name = ?", username).Get(user)
	if err != nil {
		rep.Retcode, rep.Msg = util.ThirdPackageError(err)
		return err
	}
	// 理论上应该不存在 Username 不存在的场景，除非用户在修改时，用户名被删除。
	if !res {
		rep.Retcode = conf.StatusUpdateNicknameFaild
		rep.Msg = conf.ErrMsg[conf.StatusUpdateNicknameFaild]
		return nil
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
	
	_, err = session.Where("usr_name = ?", username).Update(user)
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
		log.Println(err)
	}

	// 返回更新信息
	rep.Retcode = conf.StatusSuccess
	rep.Msg = conf.ErrMsg[conf.StatusSuccess]
	rep.Data = new(pb.UpdateNicknameReturn_Data)
	rep.Data.Nickname = user.Nickname
	return nil
}

func (s *Server) UploadPic(ctx context.Context, req *pb.UploadPicInfo, rep *pb.UpdatePicReturn) error{
	// tmp := pb.UpdatePicReturn{}
	return nil
}