package serverfunc

import (
	"context"
	"encoding/json"
	"log"

	conf "et-config/src/statusconfig"
	pb "et-protobuf3/src/gomicroapi"
	"tcpserver/src/db"
	"tcpserver/src/rediscache"

	"github.com/gomodule/redigo/redis"
)
type Server struct {
}

var retcode int32
var msg string

// tcp 服务端检查密码，demo
func (s *Server) UserLogin(ctx context.Context, req *pb.UserLoginInfo, rep *pb.LoginReturn) error {
	// rep.Retcode = conf.StatusSuccess
	// rep.Msg = conf.ErrMsg[conf.StatusSuccess]
	engine, err := db.DBConnect(db.DBname, db.Password)
	if err != nil {
		rep.Retcode, rep.Msg = ThirdPackageError(err)
		return err
	}
	defer engine.Close()
	user := new(db.User)
	_, err = engine.Where("usr_name = ?", req.Username).Get(user)
	if err != nil {
		rep.Retcode, rep.Msg = ThirdPackageError(err)
		return err
	}

	if user.Password != req.Password {
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
	conn := rediscache.RedisInit()
	user := new(db.User)
	res_json, err := redis.Bytes(conn.Do("GET", username))

	// cache miss
	if err != nil {
		log.Println("redis cache miss")
		// read DB, return engine after sync()
		engine, err := db.DBConnect(db.DBname, db.Password)
		if err != nil {
			rep.Retcode, rep.Msg = ThirdPackageError(err)
			return err
		}
		defer engine.Close()
		_, err = engine.Where("usr_name = ?", req.Username).Get(user)
		if err != nil {
			rep.Retcode, rep.Msg = ThirdPackageError(err)
			return err
		}

		// 数据 encode 准备打入 redis
		user_data_json_encode, err := json.Marshal(*user)
		if err != nil {
			log.Println(err)
		}else{
			_, err = conn.Do("set", user.Name, user_data_json_encode)
			if err != nil {
				log.Println(err)
			}
		}
	}else{
		log.Println("redis cache 命中")
		// 缓存命中，直接导入
		err = json.Unmarshal(res_json, user)
		if err != nil {
			rep.Retcode, rep.Msg = ThirdPackageError(err)
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

func (s *Server) UpdateNickname(ctx context.Context, in *pb.UpdateNicknameInfo, rep *pb.UpdateNicknameReturn) error{
	// tmp := pb.UpdateNicknameReturn{}
	// return &tmp, nil
	return nil
}

func (s *Server) UploadPic(ctx context.Context, in *pb.UploadPicInfo, rep *pb.UpdatePicReturn) error{
	// tmp := pb.UpdatePicReturn{}
	return nil
}