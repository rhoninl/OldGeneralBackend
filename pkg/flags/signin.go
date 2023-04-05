package flags

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/helper"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	"gorm.io/gorm"
)

func (s *server) SignInFlag(ctx context.Context, in *flagspb.SignInFlagRequest) (*flagspb.SignInFlagReply, error) {
	log.Println("sign in flag request", in.RequestId, in.Info.FlagId)
	txn := database.GetDB()
	return signinFlag(ctx, txn, in)
}

func (s *server) GetSignInInfo(ctx context.Context, in *flagspb.GetSignInInfoRequest) (*flagspb.GetSignInInfoReply, error) {
	log.Println("get sign in info request", in)

	var signIns model.SignIn
	err := database.GetDB().Model(&model.SignIn{}).Where("id = ?", in.SignInId).Find(&signIns).Error
	if err != nil {
		log.Println("error getting sign in info", err)
		return nil, err
	}

	signIn, err := helper.TypeConverter[cdr.SignInInfo](signIns)
	if err != nil {
		log.Println("error converting sign in info", err)
		return nil, err
	}
	reply := &flagspb.GetSignInInfoReply{
		RequestId: in.RequestId,
		ReplyTime: helper.GetTimeStamp(),
		Info:      signIn,
	}
	return reply, nil
}

func signinFlag(ctx context.Context, txn *gorm.DB, in *flagspb.SignInFlagRequest) (*flagspb.SignInFlagReply, error) {
	signinInfo, err := helper.TypeConverter[model.SignIn](in.Info)
	if err != nil {
		log.Println("error converting flag", err)
		return nil, err
	}
	var flag model.FlagInfo
	err = txn.Where("id = ?", in.Info.FlagId).Find(&flag).Error
	if err != nil {
		log.Println("error getting flag info", err)
		return nil, err
	}

	var counter int64
	err = txn.Model(&model.SignIn{}).Where("flag_id = ?", flag.ID).Where("current_time = ?", in.Info.CurrentTime).Count(&counter).Error
	if err != nil {
		log.Println("error getting sign in info", err)
		return nil, err
	}

	if counter > 0 {
		log.Println("already signed in")
		return nil, errors.New("already signed in")
	}

	if in.Info.CurrentTime == int64(flag.TotalTime) {
		flag.Status = "finished"
	}

	signinInfo.ID = in.Info.Id
	signinInfo.CreatedAt = time.Now().UnixMicro()
	signinInfo.UserID = flag.UserID
	signinInfo.TotalTime = flag.TotalTime

	err = txn.Model(&signinInfo).Save(&signinInfo).Error
	if err != nil {
		log.Println("error saving sign in info", err)
		return nil, err
	}

	var reply = &flagspb.SignInFlagReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}

	return reply, nil
}

func getSignInlist(txn *gorm.DB, flagId string) []*cdr.SignInInfo {
	var signInfos []*model.SignIn
	err := txn.Model(&model.SignIn{}).Where("flag_id = ?", flagId).Order("`current_time` DESC").Find(&signInfos).Error
	if err != nil {
		log.Println("error getting sign in info", err)
		return nil
	}
	var signList []*cdr.SignInInfo
	for _, item := range signInfos {
		flagItem, err := helper.TypeConverter[cdr.SignInInfo](item)
		if err != nil {
			log.Println("error converting sign in info", err)
			return nil
		}
		signList = append(signList, flagItem)
	}
	return signList
}
