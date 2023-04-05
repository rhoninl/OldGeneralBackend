package flags

import (
	"context"
	"log"
	"time"

	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	vippb "github.com/leepala/OldGeneralBackend/Proto/vip"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	"github.com/leepala/OldGeneralBackend/pkg/vip"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (s *server) AskForSkip(ctx context.Context, in *flagspb.AskForSkipRequest) (*flagspb.AskForSkipReply, error) {
	log.Println("ask for skip request", in)
	var prop = &model.Prop{
		ID:     uuid.NewV4().String(),
		FlagID: in.FlagId,
		Type:   int32(cdr.PropType_skip),
		UseAt:  time.Now().UnixMicro(),
	}
	err := database.GetDB().Model(prop).Save(prop).Error
	if err != nil {
		log.Println("error saving prop info", err)
		return nil, err
	}

	reply := &flagspb.AskForSkipReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}

func (s *server) Resurrect(ctx context.Context, in *flagspb.ResurrectRequest) (*flagspb.ResurrectReply, error) {
	log.Println("resurrect request", in)
	var prop = &model.Prop{
		ID:     uuid.NewV4().String(),
		FlagID: in.FlagId,
		Type:   int32(cdr.PropType_resurrection),
		UseAt:  time.Now().UnixMicro(),
	}
	err := database.GetDB().Model(prop).Save(prop).Error
	if err != nil {
		log.Println("error saving prop info", err)
		return nil, err
	}

	reply := &flagspb.ResurrectReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}

func getSkipCardUsedNum(txn *gorm.DB, flagId string) (int64, error) {
	var counter int64
	err := txn.Model(&model.Prop{}).Where("flag_id = ? and type = ?", flagId, cdr.PropType_skip).Count(&counter).Error
	if err != nil {
		log.Println("error getting skip card used num", err)
		return 0, err
	}
	return counter, nil
}

func getResurrectUsedNum(txn *gorm.DB, flagId string) (int64, error) {
	var counter int64
	err := txn.Model(&model.Prop{}).Where("flag_id = ? and type = ?", flagId, cdr.PropType_resurrection).Count(&counter).Error
	if err != nil {
		log.Println("error getting resurrect used num", err)
		return 0, err
	}
	return counter, nil
}

func dayToMaskNum(userId string, day int64) (int64, error) {
	req := &vippb.GetVipStatusRequest{
		RequestId:   uuid.NewV4().String(),
		RequestTime: time.Now().UnixMicro(),
		UserId:      userId,
	}
	resp, err := vip.GetClient().GetVipStatus(context.Background(), req)
	if err != nil {
		log.Println("error getting vip status", err)
		return 0, err
	}
	isVIP := resp.EndTime < time.Now().UnixMicro()
	maskNum := day / 7
	if isVIP {
		maskNum *= 2
	}

	return maskNum, nil
}

func dayToResurrect(day int64) int64 {
	if day >= 30 {
		return 1
	}
	return 0
}
