package api

import (
	"context"
	"log"

	vippb "github.com/leepala/OldGeneralBackend/Proto/vip"
	"github.com/leepala/OldGeneralBackend/pkg/vip"
)

func (s *server) ChargeVip(ctx context.Context, in *vippb.ChargeVipRequest) (*vippb.ChargeVipReply, error) {
	log.Println("charge vip request", in.RequestId, in.UserId, in.ChargeDuration)
	return vip.GetClient().ChargeVip(ctx, in)
}

func (s *server) GetVipStatus(ctx context.Context, in *vippb.GetVipStatusRequest) (*vippb.GetVipStatusReply, error) {
	log.Println("get vip status request", in.RequestId, in.UserId)
	return vip.GetClient().GetVipStatus(ctx, in)
}
