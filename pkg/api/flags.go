package api

import (
	"context"
	"log"

	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	"github.com/leepala/OldGeneralBackend/pkg/flags"
)

func (s *server) SearchMyFlag(ctx context.Context, in *flagspb.SearchMyFlagRequest) (*flagspb.SearchMyFlagReply, error) {
	log.Println("get my flag request", in.RequestId, in.UserId)
	return flags.GetClient().SearchMyFlag(ctx, in)
}

func (s *server) GetFlagDetail(ctx context.Context, in *flagspb.GetFlagDetailRequest) (*flagspb.GetFlagDetailReply, error) {
	log.Println("get flag info request", in)
	return flags.GetClient().GetFlagDetail(ctx, in)
}

func (s *server) CreateFlag(ctx context.Context, in *flagspb.CreateFlagRequest) (*flagspb.CreateFlagReply, error) {
	log.Println("create flag request", in.RequestId, in.Info.Id)
	return flags.GetClient().CreateFlag(ctx, in)
}

func (s *server) FetchFlagSquare(ctx context.Context, in *flagspb.FetchFlagSquareRequest) (*flagspb.FetchFlagSquareReply, error) {
	log.Println("fetch flag square request", in)
	return flags.GetClient().FetchFlagSquare(ctx, in)
}

func (s *server) SearchFlag(ctx context.Context, in *flagspb.SearchFlagRequest) (*flagspb.SearchFlagReply, error) {
	log.Println("search flag request", in)
	return flags.GetClient().SearchFlag(ctx, in)
}

func (s *server) GetSignInInfo(ctx context.Context, in *flagspb.GetSignInInfoRequest) (*flagspb.GetSignInInfoReply, error) {
	log.Println("get sign in info request", in)
	return flags.GetClient().GetSignInInfo(ctx, in)
}

func (s *server) SignInFlag(ctx context.Context, in *flagspb.SignInFlagRequest) (*flagspb.SignInFlagReply, error) {
	log.Println("sign in flag request", in)
	return flags.GetClient().SignInFlag(ctx, in)
}

func (s *server) SiegeFlag(ctx context.Context, in *flagspb.SiegeFlagRequest) (*flagspb.SiegeFlagReply, error) {
	log.Println("siege flag request", in)
	return flags.GetClient().SiegeFlag(ctx, in)
}

func (s *server) GetMySiegeNum(ctx context.Context, in *flagspb.GetMySiegeNumRequest) (*flagspb.GetMySiegeNumReply, error) {
	log.Println("get my siege num request", in)
	return flags.GetClient().GetMySiegeNum(ctx, in)
}

func (s *server) FetchMySiege(ctx context.Context, in *flagspb.FetchMySiegeRequest) (*flagspb.FetchMySiegeReply, error) {
	log.Println("fetch my siege request", in)
	return flags.GetClient().FetchMySiege(ctx, in)
}

func (s *server) CheckIsSiege(ctx context.Context, in *flagspb.CheckIsSiegeRequest) (*flagspb.CheckIsSiegeReply, error) {
	log.Println("check is sieged request", in)
	return flags.GetClient().CheckIsSiege(ctx, in)
}

func (s *server) PostComment(ctx context.Context, in *flagspb.PostCommentRequest) (*flagspb.PostCommentReply, error) {
	log.Println("post comment request", in)
	return flags.GetClient().PostComment(ctx, in)
}

func (s *server) FetchComment(ctx context.Context, in *flagspb.FetchCommentRequest) (*flagspb.FetchCommentReply, error) {
	log.Println("fetch comment request", in)
	return flags.GetClient().FetchComment(ctx, in)
}
