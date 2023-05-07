package iam

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/leepala/OldGeneralBackend/Proto/iam"
	iampb "github.com/leepala/OldGeneralBackend/Proto/iam"
	userpb "github.com/leepala/OldGeneralBackend/Proto/user"
	walletpb "github.com/leepala/OldGeneralBackend/Proto/wallet"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/helper"
	"github.com/leepala/OldGeneralBackend/pkg/helper/email"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	"github.com/leepala/OldGeneralBackend/pkg/user"
	"github.com/leepala/OldGeneralBackend/pkg/wallet"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

const (
	listenPort = ":30001"
)

type server struct {
	iampb.UnimplementedIamServer
}

func StartAndListen() {
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	iampb.RegisterIamServer(s, &server{})
	log.Println("API Server is listening on port", listenPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) IAMLogin(ctx context.Context, in *iampb.IamLoginRequest) (*iampb.IamLoginReply, error) {
	log.Println("login request", in.RequestId, in.UserName)
	var err error
	var user *model.User
	if in.Password != "" {
		user, err = loginByPassword(ctx, in.UserName, in.Password)
	} else if in.VerificationCode != "" {
		user, err = loginByEmail(ctx, in.UserName, in.VerificationCode)
	}

	if err != nil {
		return nil, err
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

func loginByPassword(ctx context.Context, username string, password string) (*model.User, error) {
	var user = &model.User{}
	err := database.GetDB().Model(&user).Where("username = ? and password = ?", username, password).Find(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == "" {
		return nil, errors.New("Cannot find user")
	}
	return user, nil
}

func loginByEmail(ctx context.Context, address string, verificationCode string) (*model.User, error) {
	cmd := database.GetRDB().GetDel(ctx, address)
	if cmd.Val() != verificationCode {
		return nil, errors.New("wrong verificationCode")
	}

	var user = &model.User{}
	var userInfo = &model.UserInfo{}
	err := database.GetDB().Model(&userInfo).Where("mail = ?").Find(&userInfo).Error
	if err != nil {
		return nil, err
	}
	if user.ID == "" {
		return nil, errors.New("User not find")
	}
	err = database.GetDB().Model(&user).Where("id = ?").Find(&user).Error
	return user, err
}

func (s *server) IAMRegister(ctx context.Context, in *iampb.CreateUserRequest) (*iampb.CreateUserReply, error) {
	log.Println("regist request", in.RequestId, in.UserName)
	var counter int64 = 1
	var userInfo = &model.User{
		ID:       uuid.NewV4().String(),
		Username: in.UserName,
		Password: in.Password,
	}
	var reply = &iampb.CreateUserReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		IsSuccess: false,
	}

	err := database.GetDB().Model(&userInfo).Where("username = ?", in.UserName).Count(&counter).Error
	if err != nil {
		log.Println("Error to get user from db", err)
		return nil, err
	}

	if counter > 0 {
		var msg = "User already exists"
		reply.Message = &msg
		return reply, nil
	}

	err = database.GetDB().Model(&userInfo).Save(&userInfo).Error
	if err != nil {
		log.Println("Error to save user to db", err)
		return nil, err
	}

	initRequest := &userpb.InitUserInfoRequest{
		UserId: userInfo.ID,
	}

	_, err = user.GetClient().InitUserInfo(ctx, initRequest)
	if err != nil {
		log.Println("Error to init user info", err)
		return nil, err
	}

	initWalletRequest := &walletpb.InitWalletRequest{
		RequestId:   in.RequestId,
		RequestTime: in.RequestTime,
		UserId:      userInfo.ID,
	}
	_, err = wallet.GetClient().InitWallet(ctx, initWalletRequest)
	if err != nil {
		log.Println("Error to init wallet", err)
		return nil, err
	}

	reply.IsSuccess = true
	return reply, nil
}

func (s *server) IAMCheckLoginStatus(ctx context.Context, in *iampb.IamCheckStatusRequest) (*iampb.IamCheckStatusReply, error) {
	log.Println("check login status request", in.RequestId)

	userId, err := helper.GetUserIdFromContext(ctx)
	if err != nil {
		log.Println("Error to get user id from context", err)
		return nil, err
	}

	var user = &model.User{}
	err = database.GetDB().Model(user).Where("id = ?", userId).Find(&user).Error
	if err != nil {
		log.Println("Error to get user from db", err)
		return nil, err
	}

	token, err := helper.GenerateToken(user.ID)
	if err != nil {
		log.Println("Error to generate token", err)
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

func (s *server) IAMSendMail(ctx context.Context, in *iam.SendMailRequest) (*iam.SendMailReply, error) {
	log.Println("send mail request, address: ", in.Address)
	err := sendVerificationCode(ctx, in.Address)
	if err != nil {
		log.Println("error to send verification code, error: ", err.Error())
		return nil, err
	}

	return &iam.SendMailReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}, nil
}

func sendVerificationCode(ctx context.Context, address string) error {
	verificationCode := email.GenerateVerificationCode()
	err := email.SendCode(verificationCode, address)
	if err != nil {
		log.Println("error to send mail, error: ", err.Error())
		return err
	}

	cmd := database.GetRDB().Set(ctx, address, verificationCode, time.Hour)
	return cmd.Err()
}
