package flags

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	"github.com/leepala/OldGeneralBackend/Proto/flags"
	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	vippb "github.com/leepala/OldGeneralBackend/Proto/vip"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/helper"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	"github.com/leepala/OldGeneralBackend/pkg/vip"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (s *server) AskForSkip(ctx context.Context, in *flagspb.AskForSkipRequest) (reply *flagspb.AskForSkipReply, err error) {
	log.Println("ask for skip request", in)
	txn := database.GetDB().Begin()
	defer helper.TransactionHandle(txn, &err)
	return askForSkip(ctx, txn, in)
}

func (s *server) Resurrect(ctx context.Context, in *flagspb.ResurrectRequest) (reply *flagspb.ResurrectReply, err error) {
	log.Println("resurrect request", in)
	txn := database.GetDB().Begin()
	defer helper.TransactionHandle(txn, &err)
	return resurrectFlag(ctx, txn, in)
}

func (s *server) WaiverResurrect(ctx context.Context, in *flagspb.WaiverResurrectRequest) (reply *flagspb.WaiverResurrectReply, err error) {
	log.Println("waiver resurrect request", in)
	txn := database.GetDB().Begin()
	defer helper.TransactionHandle(txn, &err)
	return waiverResurrect(ctx, txn, in)
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

func askForSkip(ctx context.Context, txn *gorm.DB, in *flags.AskForSkipRequest) (*flags.AskForSkipReply, error) {
	var prop = &model.Prop{
		ID:     helper.GenerateUUID(),
		FlagID: in.FlagId,
		Type:   int32(cdr.PropType_skip),
		UseAt:  helper.GetTimeStamp(),
	}
	err := txn.Model(prop).Save(prop).Error
	if err != nil {
		log.Println("error saving prop info", err)
		return nil, err
	}
	signinList := getSignInlist(txn, in.FlagId)
	signinRequest := &flagspb.SignInFlagRequest{
		RequestId:   in.RequestId,
		RequestTime: in.RequestTime,
		Info: &cdr.SignInInfo{
			Id:          helper.GenerateUUID(),
			FlagId:      in.FlagId,
			PictureUrl:  holidayUrl,
			Content:     "今天放假，不打卡",
			CurrentTime: int64(len(signinList)) + 1,
			CreatedAt:   helper.GetTimeStamp(),
			IsSkip:      1,
		},
	}
	_, err = signinFlag(ctx, txn, signinRequest)
	if err != nil {
		log.Println("error signing in", err)
		return nil, err
	}
	reply := &flagspb.AskForSkipReply{
		RequestId: in.RequestId,
		ReplyTime: helper.GetTimeStamp(),
	}
	return reply, nil
}

func updateStatusToResurrect(txn *gorm.DB, info *model.FlagInfo) error {
	var flagInfo model.FlagInfo
	err := txn.Model(&flagInfo).Where("id = ?", info.ID).Find(&flagInfo).Error
	if err != nil {
		log.Println("error getting flag info", err)
		return err
	}
	if info.Status != "running" {
		log.Println("flag is not running")
		return errors.New("flag is not running")
	}
	flagInfo.Status = "resurrect"
	err = txn.Model(&flagInfo).Where("id = ?", info.ID).Save(&flagInfo).Error
	return err
}

func resurrectFlag(ctx context.Context, txn *gorm.DB, in *flags.ResurrectRequest) (*flags.ResurrectReply, error) {
	log.Println("resurrect request", in)

	// get flagInfo
	var flagInfo *model.FlagInfo
	err := txn.Model(&model.FlagInfo{}).Where("id = ?", in.FlagId).Find(&flagInfo).Error
	if err != nil {
		log.Println("error getting flag info", err)
		return nil, err
	}
	lastSignInTime, need := needSignInToday(txn, flagInfo)
	if !need {
		log.Println("flag cannot resurrect")
		return nil, errors.New("flag cannot resurrect")
	}

	// check if user has resurrect card
	resurrectUsedNum, err := getResurrectUsedNum(txn, in.FlagId)
	if err != nil {
		log.Println("error getting resurrect used num", err)
		return nil, err
	}

	if resurrectUsedNum >= int64(flagInfo.TotalResurrectNum) {
		log.Println("resurrect card used up")
		return nil, errors.New("resurrect card used up")
	}

	// try to signin
	resurrectSignin := &cdr.SignInInfo{
		Id:          helper.GenerateUUID(),
		FlagId:      in.FlagId,
		CreatedAt:   helper.GetTimeStamp(),
		CurrentTime: lastSignInTime + 1,
		IsSkip:      2,
		Content:     "使用了复活卡",
		PictureUrl:  resurrectUrl,
	}
	signinRequest := &flagspb.SignInFlagRequest{
		RequestId:   in.RequestId,
		RequestTime: in.RequestTime,
		Info:        resurrectSignin,
	}

	_, err = signinFlag(ctx, txn, signinRequest)
	if err != nil {
		log.Println("error signing in", err)
		return nil, err
	}

	// add record of resurrect card
	var prop = &model.Prop{
		ID:     uuid.NewV4().String(),
		FlagID: in.FlagId,
		Type:   int32(cdr.PropType_resurrection),
		UseAt:  time.Now().UnixMicro(),
	}

	err = txn.Model(prop).Save(prop).Error
	if err != nil {
		log.Println("error saving prop info", err)
		return nil, err
	}

	err = txn.Model(&model.FlagInfo{}).Where("id = ?", in.FlagId).Update("status", "running").Error
	if err != nil {
		log.Println("error updating flag status", err)
		return nil, err
	}

	reply := &flagspb.ResurrectReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}

func waiverResurrect(ctx context.Context, txn *gorm.DB, in *flagspb.WaiverResurrectRequest) (*flagspb.WaiverResurrectReply, error) {
	log.Println("waiver resurrect request", in)

	// get flagInfo
	var flagInfo *model.FlagInfo
	err := txn.Model(&model.FlagInfo{}).Where("id = ?", in.FlagId).Find(&flagInfo).Error
	if err != nil {
		log.Println("error getting flag info", err)
		return nil, err
	}
	if flagInfo.Status != "resurrect" {
		log.Println("flag is not in resurrect status")
		return nil, errors.New("flag is not in resurrect status")
	}

	err = txn.Model(&model.FlagInfo{}).Where("id = ?", in.FlagId).Update("status", "failed").Error
	if err != nil {
		log.Println("error updating flag status", err)
		return nil, err
	}
	return nil, err
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
	isVIP := resp.EndTime > time.Now().UnixMicro()
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
