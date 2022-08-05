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
		{&pb.UserLoginInfo{Username: "Ding", Password: "dingyuan1"}, &pb.LoginReturn{}, &pb.LoginReturn{Retcode: conf.StatusLoginFailed}},
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