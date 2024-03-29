package api

import (
	"context"
	"log"

	walletpb "github.com/leepala/OldGeneralBackend/Proto/wallet"
	"github.com/leepala/OldGeneralBackend/pkg/wallet"
)

func (s *server) GetCurrentGold(ctx context.Context, in *walletpb.GetCurrentGoldRequest) (*walletpb.GetCurrentGoldReply, error) {
	log.Println("get current money request", in.RequestId, in.UserId)
	return wallet.GetClient().GetCurrentGold(ctx, in)
}

func (s *server) UpdateGold(ctx context.Context, in *walletpb.UpdateGoldRequest) (*walletpb.UpdateGoldReply, error) {
	log.Println("update gold request", in.RequestId, in.UserId)
	return wallet.GetClient().UpdateGold(ctx, in)
}

func (s *server) FetchWaterFlow(ctx context.Context, in *walletpb.FetchWaterFlowRequest) (*walletpb.FetchWaterFlowReply, error) {
	log.Println("fetch water flow request", in.RequestId, in.UserId)
	return wallet.GetClient().FetchWaterFlow(ctx, in)
}
