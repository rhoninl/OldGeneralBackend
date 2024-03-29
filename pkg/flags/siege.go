package flags

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	walletpb "github.com/leepala/OldGeneralBackend/Proto/wallet"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/helper"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	"github.com/leepala/OldGeneralBackend/pkg/wallet"
	uuid "github.com/satori/go.uuid"
)

func (s *server) SiegeFlag(ctx context.Context, in *flagspb.SiegeFlagRequest) (*flagspb.SiegeFlagReply, error) {
	log.Println("siege flag request", in)

	txn := database.GetDB()

	// check flag status
	var flag model.FlagInfo
	err := txn.Model(&flag).Where("id = ?", in.FlagId).Find(&flag).Error
	if err != nil {
		log.Println("error getting flag info", err)
		return nil, err
	}

	if flag.Status != "running" {
		log.Println("flag not running")
		return nil, errors.New("flag not running")
	}

	// check have sieged
	var counter int64
	err = txn.Model(&model.Siege{}).Where("flag_id = ? and user_id = ?", in.FlagId, in.UserId).Count(&counter).Error
	if err != nil {
		log.Println("error getting siege info", err)
		return nil, err
	}

	if counter > 0 {
		log.Println("user already siege")
		return nil, errors.New("user already siege")
	}

	// update money
	updateGoldRequest := &walletpb.UpdateGoldRequest{
		RequestId:   in.RequestId,
		RequestTime: in.RequestTime,
		UserId:      in.UserId,
		GoldNum:     -10,
		Content:     "围观Flag: " + flag.Name,
	}
	_, err = wallet.GetClient().UpdateGold(ctx, updateGoldRequest)
	if err != nil {
		log.Println("error updating gold", err)
		return nil, err
	}

	// try to siege
	siege := model.Siege{
		ID:        uuid.NewV4().String(),
		FlagID:    in.FlagId,
		UserID:    in.UserId,
		CreatedAt: time.Now().UnixMicro(),
	}
	err = database.GetDB().Model(&siege).Save(&siege).Error
	if err != nil {
		log.Println("error saving siege info", err)
		return nil, err
	}

	reply := &flagspb.SiegeFlagReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}

func (s *server) CheckIsSiege(ctx context.Context, in *flagspb.CheckIsSiegeRequest) (*flagspb.CheckIsSiegeReply, error) {
	log.Println("check is siege request", in)

	var siege model.Siege
	err := database.GetDB().Model(&model.Siege{}).Where("flag_id = ? and user_id = ?", in.FlagId, in.UserId).Find(&siege).Error
	if err != nil {
		log.Println("error getting siege info", err)
		return nil, err
	}

	reply := &flagspb.CheckIsSiegeReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		IsSiege:   siege.ID != "",
	}
	return reply, nil
}

func (s *server) GetMySiegeNum(ctx context.Context, in *flagspb.GetMySiegeNumRequest) (*flagspb.GetMySiegeNumReply, error) {
	log.Println("get my siege num request", in)

	var counter int64
	err := database.GetDB().Model(&model.Siege{}).Where("user_id = ?", in.UserId).Count(&counter).Error
	if err != nil {
		log.Println("error getting siege info", err)
		return nil, err
	}

	reply := &flagspb.GetMySiegeNumReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		Num:       counter,
	}
	return reply, nil
}

func (s *server) FetchMySiege(ctx context.Context, in *flagspb.FetchMySiegeRequest) (*flagspb.FetchMySiegeReply, error) {
	log.Println("fetch my siege request", in)

	var lastSiegeTimeStamp int64 = time.Now().UnixMicro() + 1
	if in.LastSiegeId != "" {
		var siege model.Siege
		err := database.GetDB().Model(&model.Siege{}).Where("id = ?", in.LastSiegeId).Find(&siege).Error
		if err != nil {
			log.Println("error getting siege info", err)
			return nil, err
		}
		lastSiegeTimeStamp = siege.CreatedAt
	}
	var sieges []model.Siege
	err := database.GetDB().Model(&model.Siege{}).Where("user_id = ? and created_at < ?", in.UserId, lastSiegeTimeStamp).Limit(int(in.PageSize)).Find(&sieges).Error
	if err != nil {
		log.Println("error getting siege info", err)
		return nil, err
	}

	var siegeInfos []*cdr.SiegeInfo
	var tmpsiege *cdr.SiegeInfo
	for _, siege := range sieges {
		tmpsiege, err = helper.TypeConverter[cdr.SiegeInfo](siege)
		if err != nil {
			log.Println("error converting siege info", err)
			return nil, err
		}
		siegeInfos = append(siegeInfos, tmpsiege)
	}

	reply := &flagspb.FetchMySiegeReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		Sieges:    siegeInfos,
	}
	return reply, nil
}

func getSiegeNumByFlagId(flagId string) int64 {
	var counter int64
	err := database.GetDB().Model(&model.Siege{}).Where("flag_id = ?", flagId).Count(&counter).Error
	if err != nil {
		log.Println("error getting siege info", err)
		return 0
	}
	return counter
}
