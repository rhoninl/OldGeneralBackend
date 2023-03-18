package user

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	userpb "github.com/leepala/OldGeneralBackend/Proto/user"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/helper"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	listenPort    = ":30001"
	defaultAvatar = "https://oldgeneral.obs.cn-north-4.myhuaweicloud.com/avatars/turtlerock.jpg"
)

type server struct {
	userpb.UnimplementedUserServer
}

func StartAndListen() {
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	userpb.RegisterUserServer(s, &server{})
	log.Println("API Server is listening on port", listenPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) GetUserInfo(ctx context.Context, in *userpb.GetUserInfoRequest) (*userpb.GetUserInfoReply, error) {
	log.Println("login request", in.RequestId, in.UserId)
	var user = &model.UserInfo{}
	err := database.GetDB().Model(&user).Where("user_id = ?", in.UserId).Find(&user).Error
	if err != nil {
		log.Printf("cannot get userBasic Info by user id: %s, error: %s", in.UserId, err)
		return nil, errors.New("Cannot find user")
	}
	userInfoReply, err := helper.TypeConverter[cdr.UserBasicInfo](user)
	if err != nil {
		log.Printf("cannot convert userBasic Info to reply, error: %s", err)
		return nil, errors.New("userInfo is invalid")
	}
	var reply = &userpb.GetUserInfoReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		UserInfo:  userInfoReply,
	}
	return reply, nil
}

func (s *server) UpdateUserInfo(ctx context.Context, in *userpb.UpdateUserInfoRequest) (*userpb.UpdateUserInfoReply, error) {
	log.Println("update user info request")
	userId, err := helper.GetUserIdFromContext(ctx)
	if err != nil {
		log.Printf("cannot get user id from context, error: %s", err)
		return nil, errors.New("cannot get user id from context")
	}
	userInfo := &model.UserInfo{}
	err = database.GetDB().Model(&model.UserInfo{}).Where("user_id = ?", userId).Find(&userInfo).Error
	if err != nil {
		log.Printf("cannot get userBasic Info by user id: %s, error: %s", userId, err)
		return nil, errors.New("Cannot find user")
	}

	if in.UserName != nil {
		userInfo.Name = *in.UserName
	}
	if in.UserAvatar != nil {
		userInfo.Avatar = *in.UserAvatar
	}
	if in.UserGender != nil {
		userInfo.Gender = *in.UserGender
	}
	if in.UserSignature != nil {
		userInfo.Signature = *in.UserSignature
	}

	err = database.GetDB().Model(&model.UserInfo{}).Where("user_id = ?", userId).Save(&userInfo).Error
	if err != nil {
		log.Printf("cannot update userBasic Info by user id: %s, error: %s", userId, err)
		return nil, errors.New("Cannot update user")
	}
	var reply = &userpb.UpdateUserInfoReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}

func (s *server) InitUserInfo(ctx context.Context, in *userpb.InitUserInfoRequest) (*emptypb.Empty, error) {
	log.Println("init user info request userId ", in.UserId)
	userInfo := &model.UserInfo{
		ID:        in.UserId,
		Name:      "将军",
		Gender:    "男",
		CreatedAt: time.Now(),
		Avatar:    defaultAvatar,
		Signature: "",
	}
	err := database.GetDB().Model(&model.UserInfo{}).Create(&userInfo).Error
	if err != nil {
		log.Printf("cannot init userBasic Info by user id: %s, error: %s", in.UserId, err)
		return nil, errors.New("Cannot init user")
	}
	return &emptypb.Empty{}, nil
}
