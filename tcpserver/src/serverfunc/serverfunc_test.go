package serverfunc

import (
	"context"
	pb "et-protobuf3/src/gomicroapi"
	"testing"

	conf "et-config/src/statusconfig"
)

func TestUserLogin(t *testing.T) {
	tests := []struct{
		req *pb.UserLoginInfo
		rep *pb.LoginReturn
		want *pb.LoginReturn
	}{
		{&pb.UserLoginInfo{Username: "Ding", Password: "dingyuan"}, &pb.LoginReturn{}, &pb.LoginReturn{Retcode: conf.StatusSuccess}},
		{&pb.UserLoginInfo{Username: "Ding", Password: "dingyuan1"}, &pb.LoginReturn{}, &pb.LoginReturn{Retcode: conf.StatusLoginFailed}}, // password wrong
		{&pb.UserLoginInfo{Username: "Ding1", Password: "dingyuan"}, &pb.LoginReturn{}, &pb.LoginReturn{Retcode: conf.StatusLoginFailed}}, // user not exits
	}
	ctx := context.Background()
	s := Server{}
	for _, test := range tests {
		err := s.UserLogin(ctx, test.req, test.rep)
		if err != nil {
			t.Error(err)
		}
		if test.rep.Retcode != test.want.Retcode {
			t.Errorf("reture code err: got %d, wand %d", test.rep.Retcode, test.want.Retcode)
		}
	}
}

func TestUserQuery(t *testing.T) {
	tests := []struct{
		req *pb.UserQueryInfo
		rep *pb.QueryReturn
		want *pb.QueryReturn
	}{
		{&pb.UserQueryInfo{Username: "Ding"}, &pb.QueryReturn{}, &pb.QueryReturn{Retcode: conf.StatusSuccess}},
		{&pb.UserQueryInfo{Username: "Ding1"}, &pb.QueryReturn{}, &pb.QueryReturn{Retcode: conf.StatusQueryFaild}}, // user not exits
	}
	ctx := context.Background()
	s := Server{}
	for _, test := range tests {
		err := s.UserQuery(ctx, test.req, test.rep)
		if err != nil {
			t.Error(err)
		}
		if test.rep.Retcode != test.want.Retcode {
			t.Errorf("reture code err: got %d, wand %d", test.rep.Retcode, test.want.Retcode)
		}
	}
}

func TestUpdateNickname(t *testing.T) {
	tests := []struct{
		req *pb.UpdateNicknameInfo
		rep *pb.UpdateNicknameReturn
		want *pb.UpdateNicknameReturn
	}{
		// {&pb.UpdateNicknameInfo{Username: "Ding", Nickname: "dingyuan123"}, &pb.UpdateNicknameReturn{}, &pb.UpdateNicknameReturn{Retcode: conf.StatusSuccess}},
		// {&pb.UpdateNicknameInfo{Username: "Ding", Nickname: "dingyuan231231231231231232132131231231231231231231231231231231231231231"}, &pb.UpdateNicknameReturn{}, &pb.UpdateNicknameReturn{Retcode: conf.StatusUpdateNicknameFaild}}, // nickname to long
		{&pb.UpdateNicknameInfo{Username: "Ding1", Nickname: "dingyuan"}, &pb.UpdateNicknameReturn{}, &pb.UpdateNicknameReturn{Retcode: conf.StatusUpdateNicknameFaild}}, // user not exits
	}

	ctx := context.Background()
	s := Server{}
	for _, test := range tests {
		err := s.UpdateNickname(ctx, test.req, test.rep)
		if err != nil {
			t.Error(err)
		}
		if test.rep.Retcode != test.want.Retcode {
			t.Errorf("reture code err: got %d, wand %d", test.rep.Retcode, test.want.Retcode)
		}
	}
}