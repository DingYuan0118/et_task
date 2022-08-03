package serverfunc

import (
	"context"
	"log"

	conf "et-config/src/statusconfig"
	pb "et-protobuf3/src/gomicroapi"
	"tcpserver/src/db"

)
type Server struct {
}

// tcp 服务端检查密码，demo
func (s *Server) UserLogin(ctx context.Context, req *pb.UserLoginInfo, rep *pb.LoginReturn) error {
	var retcode int32
	var msg string
	// rep.Retcode = conf.StatusSuccess
	// rep.Msg = conf.ErrMsg[conf.StatusSuccess]
	engine, err := db.DBConnect(db.DBname, db.Password)
	if err != nil {
		log.Println(err.Error())
		retcode = conf.StatusThirdPackageErr
		msg = err.Error()
		rep.Retcode = retcode
		rep.Msg = msg
		return err
	}
	defer engine.Close()
	user := new(db.User)
	_, err = engine.Where("usr_name = ?", req.Username).Get(user)
	if err != nil {
		log.Println(err.Error())
		retcode = conf.StatusThirdPackageErr
		msg = err.Error()
		rep.Retcode = retcode
		rep.Msg = msg
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