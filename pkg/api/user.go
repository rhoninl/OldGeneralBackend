package api

import (
	"context"
	"log"

	userpb "github.com/leepala/OldGeneralBackend/Proto/user"
	"github.com/leepala/OldGeneralBackend/pkg/user"
)

func (s *server) GetUserInfo(ctx context.Context, in *userpb.GetUserInfoRequest) (*userpb.GetUserInfoReply, error) {
	log.Println("get user info request", in.RequestId, in.UserId)
	return user.GetClient().GetUserInfo(ctx, in)
}

func (s *server) UpdateUserInfo(ctx context.Context, in *userpb.UpdateUserInfoRequest) (*userpb.UpdateUserInfoReply, error) {
	log.Println("update user info request", in.RequestId)
	return user.GetClient().UpdateUserInfo(ctx, in)
}
