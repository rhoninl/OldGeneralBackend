package wallet

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	walletpb "github.com/leepala/OldGeneralBackend/Proto/wallet"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	uuid "github.com/satori/go.uuid"

	"google.golang.org/grpc"
)

const (
	listenPort = ":30001"
)

type server struct {
	walletpb.UnimplementedWalletServer
}

func StartAndListen() {
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// walletpb.RegisterFlagsServer(s, &server{})
	walletpb.RegisterWalletServer(s, &server{})
	log.Println("Wallet Server is listening on port", listenPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) InitWallet(ctx context.Context, in *walletpb.InitWalletRequest) (*walletpb.InitWalletReply, error) {
	log.Println("init wallet request", in.RequestId, in.UserId)
	var wallet = &model.Wallet{
		ID:      uuid.NewV4().String(),
		UserID:  in.UserId,
		GoldNum: 0,
	}
	err := database.GetDB().Model(&wallet).Create(&wallet).Error
	if err != nil {
		log.Println("error creating wallet", err)
		return nil, err
	}
	var reply = &walletpb.InitWalletReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}

func (s *server) GetCurrentGold(ctx context.Context, in *walletpb.GetCurrentGoldRequest) (*walletpb.GetCurrentGoldReply, error) {
	log.Println("get current gold request", in.RequestId, in.UserId)
	var wallet = &model.Wallet{}
	err := database.GetDB().Model(&wallet).Where("user_id = ?", in.UserId).First(&wallet).Error
	if err != nil {
		log.Println("error getting wallet", err)
		return nil, err
	}
	var reply = &walletpb.GetCurrentGoldReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		GoldNum:   wallet.GoldNum,
	}
	return reply, nil
}

func (s *server) UpdateGold(ctx context.Context, in *walletpb.UpdateGoldRequest) (*walletpb.UpdateGoldReply, error) {
	log.Println("update gold request", in.RequestId, in.UserId, in.GoldNum)
	var wallet = &model.Wallet{}
	err := database.GetDB().Model(&wallet).Where("user_id = ?", in.UserId).First(&wallet).Error
	if err != nil {
		log.Println("error getting wallet", err)
		return nil, err
	}
	wallet.GoldNum += in.GoldNum
	if wallet.GoldNum < 0 {
		return nil, fmt.Errorf("not enough gold, old money: %d", wallet.GoldNum)
	}
	err = database.GetDB().Model(&wallet).Save(&wallet).Error
	if err != nil {
		log.Println("error updating wallet", err)
		return nil, err
	}
	var reply = &walletpb.UpdateGoldReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}
