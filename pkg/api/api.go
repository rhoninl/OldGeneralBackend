package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	apipb "github.com/leepala/OldGeneralBackend/Proto/api"
	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	"github.com/leepala/OldGeneralBackend/Proto/flags"
	"github.com/leepala/OldGeneralBackend/Proto/iam"
	iampb "github.com/leepala/OldGeneralBackend/Proto/iam"
	"github.com/leepala/OldGeneralBackend/Proto/userinfo"
	uuid "github.com/satori/go.uuid"

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

func (s *server) SearchMyFlag(ctx context.Context, in *flags.SearchMyFlagRequest) (*flags.SearchMyFlagReply, error) {
	log.Println("get my flag request", in.RequestId, in.UserId)
	reply := &flags.SearchMyFlagReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().Unix(),
		Flags: []*cdr.FlagBasicInfo{
			&cdr.FlagBasicInfo{
				FlagId:      uuid.NewV4().String(),
				UserId:      in.UserId,
				FlagName:    "flag1",
				FlagStatus:  "running",
				TotalTime:   100,
				CurrentTime: 50,
				StartTime:   time.Now().Unix(),
			},
			&cdr.FlagBasicInfo{
				FlagId:      uuid.NewV4().String(),
				UserId:      in.UserId,
				FlagName:    "flag2",
				FlagStatus:  "banned",
				TotalTime:   100,
				CurrentTime: 50,
				StartTime:   time.Now().Unix(),
			},
		},
	}
	return reply, nil
}

func (s *server) CreateFlag(ctx context.Context, in *flags.CreateFlagRequest) (*flags.CreateFlagReply, error) {
	log.Println("create flag request", in)
	reply := &flags.CreateFlagReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}

	return reply, nil
}

func (s *server) FetchFlagSquare(ctx context.Context, in *flags.FetchFlagSquareRequest) (*flags.FetchFlagSquareReply, error) {
	log.Println("fetch flag square request", in)
	reply := &flags.FetchFlagSquareReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	for i := 0; i < int(in.PageSize); i++ {
		flag := &cdr.FlagSquareItemInfo{
			SigninId:      uuid.NewV4().String(),
			UserName:      "MrLeea",
			FlagName:      "flag" + strconv.Itoa(i),
			TotalTime:     rand.Int63n(100) + 10,
			CurrentTime:   rand.Int63n(10),
			PayMoney:      rand.Int63n(1000),
			SiegeNum:      rand.Int63n(1000),
			SigninPicture: "https://img1.baidu.com/it/u=413417701,3210171500&fm=253&fmt=auto&app=138&f=JPEG?w=750&h=500",
		}

		reply.Flags = append(reply.Flags, flag)
	}
	return reply, nil
}

func (s *server) GetFlagDetail(ctx context.Context, in *flags.GetFlagDetailRequest) (*flags.GetFlagDetailReply, error) {
	log.Println("get flag info request", in)
	reply := &flags.GetFlagDetailReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().Unix(),
		Info: &cdr.FlagDetailInfo{
			FlagId:       in.FlagId,
			UserId:       "testUserId",
			FlagName:     "testFlag",
			FlagStatus:   "running",
			TotalTime:    100,
			CurrentTime:  0,
			StartTime:    time.Now().Unix(),
			ChallengeNum: 100,
			SiegeNum:     399,
			StarNum:      100,
			SignUpId:     []string{"123", "456", "789"},
		},
	}
	return reply, nil
}
