package api

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	apipb "github.com/leepala/OldGeneralBackend/Proto/api"
	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	"github.com/leepala/OldGeneralBackend/Proto/flags"
	uuid "github.com/satori/go.uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	apipb.UnimplementedApiServer
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, exists := metadata.FromIncomingContext(ctx)
	if exists {
		token := md.Get(CONTEXT_USER_TOKEN_AUTHORIZATION_STR)
		if len(token) > 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, CONTEXT_USER_TOKEN_AUTHORIZATION_STR, token[0])
			log.Println("api without any token")
		}
	} else {
		log.Println("cannot get authorization")
	}
	return handler(ctx, req)
}

func StartAndListen() {
	var listenPort = os.Getenv("ListenPort")
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	interceptor := grpc.UnaryInterceptor(unaryInterceptor)
	s := grpc.NewServer(interceptor)
	apipb.RegisterApiServer(s, &server{})
	log.Println("API Server is listening on port 30001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) SearchMyFlag(ctx context.Context, in *flags.SearchMyFlagRequest) (*flags.SearchMyFlagReply, error) {
	log.Println("get my flag request", in.RequestId, in.UserId)
	reply := &flags.SearchMyFlagReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		Flags: []*cdr.FlagBasicInfo{
			{
				Id:          uuid.NewV4().String(),
				UserId:      in.UserId,
				Name:        "flag1",
				Status:      "running",
				TotalTime:   100,
				CurrentTime: 50,
				StartTime:   time.Now().UnixMicro(),
			},
			{
				Id:          uuid.NewV4().String(),
				UserId:      in.UserId,
				Name:        "flag2",
				Status:      "banned",
				TotalTime:   100,
				CurrentTime: 50,
				StartTime:   time.Now().UnixMicro(),
			},
			{
				Id:          uuid.NewV4().String(),
				UserId:      in.UserId,
				Name:        "flag3",
				Status:      "finished",
				TotalTime:   100,
				CurrentTime: 50,
				StartTime:   time.Now().UnixMicro(),
			},
		},
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
		ReplyTime: time.Now().UnixMicro(),
		Info: &cdr.FlagDetailInfo{
			Id:     in.FlagId,
			UserId: "testUserId",

			Name:         "testFlag",
			Status:       "running",
			TotalTime:    100,
			CurrentTime:  0,
			StartTime:    time.Now().UnixMicro(),
			ChallengeNum: 100,
			SiegeNum:     399,
			StarNum:      100,
			SignUpId:     []string{"123", "456", "789"},
		},
	}
	return reply, nil
}
