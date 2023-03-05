package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	apipb "github.com/leepala/OldGeneralBackend/Proto/api"
	"github.com/leepala/OldGeneralBackend/Proto/iam"
	iampb "github.com/leepala/OldGeneralBackend/Proto/iam"
	"github.com/leepala/OldGeneralBackend/Proto/userinfo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	apipb.UnimplementedApiServer
}

func StartAndListen() {
	lis, err := net.Listen("tcp", ":30001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// pb.RegisterGreeterServer(s, &apipb.ApiServer{})
	apipb.RegisterApiServer(s, &server{})
	log.Println("API Server is listening on port 30001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) IAMLogin(ctx context.Context, in *iampb.IamLoginRequest) (*iampb.IamLoginReply, error) {
	log.Println("login request", in.RequestId, in.UserName, in.Password)
	reply := &iampb.IamLoginReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().Unix(),
		Token:     "123",
		UserId:    "test",
	}
	return reply, nil
}

func (s *server) IAMRegister(ctx context.Context, in *iam.CreateUserRequest) (*iam.CreateUserReply, error) {
	log.Println("regist request", in.RequestId, in.UserName, in.Password)
	reply := &iam.CreateUserReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().Unix(),
		IsSuccess: true,
	}
	return reply, nil
}

func (s *server) IAMCheckLoginStatus(ctx context.Context, in *iam.IamCheckStatusRequest) (*iam.IamCheckStatusReply, error) {
	data, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("cannot get metadata")
	}
	token := data.Get("authorization")[0]
	log.Println("check status request", in.RequestId, token)
	reply := &iam.IamCheckStatusReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().Unix(),
		IsValid:   token != "",
		UserId:    "testUserId",
	}
	return reply, nil
}

func (s *server) UserInfoGet(ctx context.Context, in *userinfo.GetUserInfoRequest) (*userinfo.GetUserInfoReply, error) {
	log.Println("get user info request", in.RequestId, in.UserId)
	reply := &userinfo.GetUserInfoReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().Unix(),
		UserInfo: &userinfo.UserBasicInfo{
			UserId:        in.UserId,
			UserName:      "MrLeea",
			UserSignature: "不要迷恋哥，哥只是个传说",
			UserAvatar:    "https://oldgeneral.obs.cn-north-4.myhuaweicloud.com:443/avatars/turtlerock.jpg",
			UserGender:    "男",
			UserBirthday:  time.Now().UnixMicro(),
		},
	}
	return reply, nil
}
