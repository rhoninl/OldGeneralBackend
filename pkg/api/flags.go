package api

import (
	"context"
	"log"

	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	"github.com/leepala/OldGeneralBackend/pkg/flags"
)

func (s *server) CreateFlag(ctx context.Context, in *flagspb.CreateFlagRequest) (*flagspb.CreateFlagReply, error) {
	log.Println("create flag request", in.RequestId, in.Info.Id)
	return flags.GetClient().CreateFlag(ctx, in)
}

func (s *server) SearchMyFlag(ctx context.Context, in *flagspb.SearchMyFlagRequest) (*flagspb.SearchMyFlagReply, error) {
	log.Println("get my flag request", in.RequestId, in.UserId)
	return flags.GetClient().SearchMyFlag(ctx, in)
}

func (s *server) GetFlagDetail(ctx context.Context, in *flagspb.GetFlagDetailRequest) (*flagspb.GetFlagDetailReply, error) {
	log.Println("get flag info request", in)
	return flags.GetClient().GetFlagDetail(ctx, in)
}

func (s *server) SignInFlag(ctx context.Context, in *flagspb.SignInFlagRequest) (*flagspb.SignInFlagReply, error) {
	log.Println("sign in flag request", in)
	return flags.GetClient().SignInFlag(ctx, in)
}

func (s *server) FetchFlagSquare(ctx context.Context, in *flagspb.FetchFlagSquareRequest) (*flagspb.FetchFlagSquareReply, error) {
	log.Println("fetch flag square request", in)
	return flags.GetClient().FetchFlagSquare(ctx, in)
}

func (s *server) GetSignInInfo(ctx context.Context, in *flagspb.GetSignInInfoRequest) (*flagspb.GetSignInInfoReply, error) {
	log.Println("get sign in info request", in)
	return flags.GetClient().GetSignInInfo(ctx, in)
}
