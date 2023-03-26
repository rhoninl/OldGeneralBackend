package wallet

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	walletpb "github.com/leepala/OldGeneralBackend/Proto/wallet"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/helper"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

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
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		err1 := tx.Model(&wallet).Where("user_id = ?", in.UserId).First(&wallet).Error
		if wallet.GoldNum < 0 {
			return fmt.Errorf("not enough gold, old money: %d", wallet.GoldNum)
		}
		wallet.GoldNum += in.GoldNum
		err2 := database.GetDB().Model(&wallet).Save(&wallet).Error
		err3 := AddWaterFullRecord(tx, in.UserId, in.GoldNum, in.Content)
		return errors.Join(err1, err2, err3)
	})
	if err != nil {
		log.Println("error getting wallet", err)
		return nil, err
	}

	var reply = &walletpb.UpdateGoldReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}

func AddWaterFullRecord(txn *gorm.DB, userId string, goldNum int64, content string) error {
	var waterFullRecord = model.WaterFlow{
		ID:        uuid.NewV4().String(),
		UserID:    userId,
		GoldNum:   goldNum,
		Content:   content,
		CreatedAt: time.Now().UnixMicro(),
	}
	err := txn.Model(&model.WaterFlow{}).Create(&waterFullRecord).Error
	if err != nil {
		log.Printf("cannot add waterFull record by user id: %s, error: %s", userId, err)
		return errors.New("Cannot add waterFull record")
	}
	return nil
}

func (s *server) FetchWaterFlow(ctx context.Context, in *walletpb.FetchWaterFlowRequest) (*walletpb.FetchWaterFlowReply, error) {
	log.Println("fetch water flow request", in.RequestId, in.UserId)
	var waterFlowRecords []*model.WaterFlow
	err := database.GetDB().Model(&model.WaterFlow{}).Where("user_id = ?", in.UserId).Order("created_at DESC").Offset(int(in.PageNum * in.PageSize)).Limit(int(in.PageSize)).Find(&waterFlowRecords).Error
	if err != nil {
		log.Printf("cannot get waterFlow records by user id: %s, error: %s", in.UserId, err)
		return nil, errors.New("Cannot find waterFlow records")
	}
	var waterfullList []*cdr.WaterFlow
	for _, record := range waterFlowRecords {
		record1, err := helper.TypeConverter[cdr.WaterFlow](record)
		if err != nil {
			log.Println("error converting water flow record", err)
			return nil, err
		}
		waterfullList = append(waterfullList, record1)
	}
	var reply = &walletpb.FetchWaterFlowReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		WaterFlow: waterfullList,
	}
	return reply, nil
}
