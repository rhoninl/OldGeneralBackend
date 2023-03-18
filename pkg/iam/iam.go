package iam

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"time"

	iampb "github.com/leepala/OldGeneralBackend/Proto/iam"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/helper"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	iampb.UnimplementedIamServer
}

func StartAndListen() {
	var listenPort = os.Getenv("ListenPort")
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	iampb.RegisterIamServer(s, &server{})
	log.Println("API Server is listening on port 30001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) IamLogin(ctx context.Context, in *iampb.IamLoginRequest) (*iampb.IamLoginReply, error) {
	log.Println("login request", in.RequestId, in.UserName)
	var user = &model.User{}
	err := database.GetDB().Model(&user).Where("username = ? and password = ?", in.UserName, in.Password).Find(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == "" {
		return nil, errors.New("Cannot find user")
	}
	token, err := helper.GenerateToken(user.ID)
	var reply = &iampb.IamLoginReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		Token:     token,
		UserId:    user.ID,
	}

	return reply, nil
}

func (s *server) IAMRegister(ctx context.Context, in *iampb.CreateUserRequest) (*iampb.CreateUserReply, error) {
	log.Println("regist request", in.RequestId, in.UserName)
	var counter int64 = 1
	var user = &model.User{
		ID:       uuid.NewV4().String(),
		Username: in.UserName,
		Password: in.Password,
	}
	var reply = &iampb.CreateUserReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		IsSuccess: false,
	}

	err := database.GetDB().Model(&user).Where("username = ?", in.UserName).Count(&counter).Error
	if err != nil {
		return nil, err
	}

	if counter > 0 {
		var msg = "User already exists"
		reply.Message = &msg
		return reply, nil
	}

	err = database.GetDB().Model(&user).Save(&user).Error
	if err != nil {
		return nil, err
	}

	reply.IsSuccess = true
	return reply, nil
}

func (s *server) IAMCheckLoginStatus(ctx context.Context, in *iampb.IamCheckStatusRequest) (*iampb.IamCheckStatusReply, error) {
	data, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("cannot get metadata")
	}

	authorization := data.Get("authorization")
	if len(authorization) == 0 {
		return nil, errors.New("cannot get authorization")
	}

	token := data.Get("authorization")[0]
	log.Println("check status request", in.RequestId, token)

	userId, legal := helper.ValidateToken(token)
	if !legal {
		return nil, errors.New("token is not valid")
	}

	var user = &model.User{}
	err := database.GetDB().Model(user).Where("id = ?", userId).Find(&user).Error
	if err != nil {
		return nil, err
	}

	token, err = helper.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	reply := &iampb.IamCheckStatusReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		IsValid:   user.ID != "",
		UserId:    user.ID,
		Token:     token,
	}
	return reply, nil
}
