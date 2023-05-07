package api

import (
	"context"
	"log"

	iampb "github.com/leepala/OldGeneralBackend/Proto/iam"
	"github.com/leepala/OldGeneralBackend/pkg/iam"
)

const (
	CONTEXT_USER_TOKEN_AUTHORIZATION_STR = "Authorization"
)

func (s *server) IAMLogin(ctx context.Context, in *iampb.IamLoginRequest) (*iampb.IamLoginReply, error) {
	log.Println("received login request", in.RequestId, in.UserName)
	return iam.GetClient().IAMLogin(ctx, in)
}

func (s *server) IAMRegister(ctx context.Context, in *iampb.CreateUserRequest) (*iampb.CreateUserReply, error) {
	log.Println("received register request", in.RequestId, in.UserName)
	return iam.GetClient().IAMRegister(ctx, in)
}

func (s *server) IAMCheckLoginStatus(ctx context.Context, in *iampb.IamCheckStatusRequest) (*iampb.IamCheckStatusReply, error) {
	log.Println("received check login status request", in.RequestId)
	return iam.GetClient().IAMCheckLoginStatus(ctx, in)
}

func (s *server) IAMSendMail(ctx context.Context, in *iampb.SendMailRequest) (*iampb.SendMailReply, error) {
	log.Println("received send mail request", in.RequestId)
	return iam.GetClient().IAMSendMail(ctx, in)
}
