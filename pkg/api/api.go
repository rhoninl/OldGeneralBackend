package api

import (
	"context"
	"log"
	"net"
	"time"

	apipb "github.com/leepala/OldGeneralBackend/Proto/api"
	"github.com/leepala/OldGeneralBackend/Proto/iam"
	iampb "github.com/leepala/OldGeneralBackend/Proto/iam"

	"google.golang.org/grpc"
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
