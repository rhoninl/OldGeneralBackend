package flags

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	userpb "github.com/leepala/OldGeneralBackend/Proto/user"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/helper"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	"github.com/leepala/OldGeneralBackend/pkg/user"
	"github.com/robfig/cron/v3"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

const (
	listenPort   = ":30001"
	holidayUrl   = "https://oldgeneral.obs.cn-north-4.myhuaweicloud.com:443/others/holiday.jpg"
	resurrectUrl = "https://oldgeneral.obs.cn-north-4.myhuaweicloud.com:443/others/resurrect.jpeg"
)

type server struct {
	flagspb.UnimplementedFlagsServer
}

func StartAndListen() {
	startCronTabJob()

	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	flagspb.RegisterFlagsServer(s, &server{})
	log.Println("Flags Server is listening on port", listenPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) CreateFlag(ctx context.Context, in *flagspb.CreateFlagRequest) (*flagspb.CreateFlagReply, error) {
	log.Println("create flag request", in.RequestId, in.Info.Id)
	flag, err := helper.TypeConverter[model.FlagInfo](in.Info)
	if err != nil {
		log.Println("error converting flag", err)
		return nil, err
	}
	maskNum, err := dayToMaskNum(flag.UserID, int64(flag.TotalTime))
	if err != nil {
		log.Println("error converting day to mask num", err)
		return nil, err
	}
	flag.TotalResurrectNum = int32(dayToResurrect(int64(flag.TotalTime)))
	flag.TotalMaskNum = int32(maskNum)
	flag.ID = uuid.NewV4().String()
	flag.CreatedAt = helper.GetTimeStamp()
	flag.Status = "running"
	err = database.GetDB().Model(&flag).Create(&flag).Error
	if err != nil {
		return nil, err
	}

	var reply = &flagspb.CreateFlagReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}

	return reply, nil
}

func (s *server) SearchMyFlag(ctx context.Context, in *flagspb.SearchMyFlagRequest) (*flagspb.SearchMyFlagReply, error) {
	log.Println("search my flag request", in.RequestId, in.UserId)
	var flags []model.FlagInfo
	err := database.GetDB().Model(&flags).Where("user_id = ?", in.UserId).Find(&flags).Error
	if err != nil {
		return nil, err
	}
	var reply = &flagspb.SearchMyFlagReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		Flags:     make([]*cdr.FlagBasicInfo, 0),
	}
	for _, flag := range flags {
		f, err := helper.TypeConverter[cdr.FlagBasicInfo](flag)
		if err != nil {
			log.Println("error converting flag", err)
			return nil, err
		}
		err = database.GetDB().Model(&model.SignIn{}).Where("flag_id = ?", flag.ID).Count(&f.CurrentTime).Error
		if err != nil {
			log.Println("error getting sign in info", err)
			return nil, err
		}
		reply.Flags = append(reply.Flags, f)
	}
	return reply, nil
}

func (s *server) GetFlagDetail(ctx context.Context, in *flagspb.GetFlagDetailRequest) (*flagspb.GetFlagDetailReply, error) {
	log.Println("get flag detail request", in.RequestId, in.FlagId)
	var flag model.FlagInfo
	txn := database.GetDB()
	err := txn.Model(&flag).Where("id = ?", in.FlagId).Find(&flag).Error
	if err != nil {
		log.Println("error getting flag info", err)
		return nil, err
	}
	f, err := helper.TypeConverter[cdr.FlagDetailInfo](flag)
	if err != nil {
		log.Println("error converting flag", err)
		return nil, err
	}
	searchUserReq := &userpb.GetUserInfoRequest{
		RequestId:   in.RequestId,
		RequestTime: in.RequestTime,
		UserId:      flag.UserID,
	}
	userInfoReply, err := user.GetClient().GetUserInfo(ctx, searchUserReq)
	if err != nil {
		log.Println("error getting user info", err)
		return nil, err
	}
	var signInfos []*model.SignIn
	err = txn.Model(&model.SignIn{}).Where("flag_id = ?", flag.ID).Find(&signInfos).Error
	if err != nil {
		log.Println("error getting sign in info", err)
		return nil, err
	}

	f.UserAvatar = userInfoReply.UserInfo.Avatar
	f.UserName = userInfoReply.UserInfo.Name
	f.SignUpInfo = getSignInlist(txn, flag.ID)
	var err1, err2 error
	f.UsedMaskNum, err1 = getSkipCardUsedNum(txn, flag.ID)
	f.UsedResurrectNum, err2 = getResurrectUsedNum(txn, flag.ID)
	err3 := database.GetDB().Model(&model.SignIn{}).Where("flag_id = ?", flag.ID).Count(&f.CurrentTime).Error
	err = errors.Join(err1, err2, err3)
	if err != nil {
		log.Println("error getting used props info", err)
		return nil, err
	}
	f.SiegeNum = getSiegeNumByFlagId(flag.ID)

	var reply = &flagspb.GetFlagDetailReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		Info:      f,
	}
	return reply, nil
}

func (s *server) FetchFlagSquare(ctx context.Context, in *flagspb.FetchFlagSquareRequest) (*flagspb.FetchFlagSquareReply, error) {
	log.Println("fetch flag square request", in)

	var lastSignInTimeStamp int64 = time.Now().UnixMicro() + 1
	if in.LastSigninId != "" {
		var lastSignInId model.FlagInfo
		err := database.GetDB().Model(&model.FlagInfo{}).Where("id = ?", in.LastSigninId).Find(&lastSignInId).Error
		if err != nil {
			log.Println("error getting last sign in info", err)
			return nil, err
		}
		lastSignInTimeStamp = lastSignInId.CreatedAt
	}

	var signIns []model.SignIn
	err := database.GetDB().Model(&model.SignIn{}).Where("is_skip = 0").Where("created_at < ?", lastSignInTimeStamp).Order("created_at DESC").Limit(int(in.PageSize)).Find(&signIns).Error
	if err != nil {
		log.Println("error getting flag info", err)
		return nil, err
	}

	reply := &flagspb.FetchFlagSquareReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	searchFlagInfoReq := &flagspb.GetFlagDetailRequest{
		RequestId:   in.RequestId,
		RequestTime: in.RequestTime,
	}
	searchUserReq := &userpb.GetUserInfoRequest{
		RequestId:   in.RequestId,
		RequestTime: in.RequestTime,
	}
	for _, item := range signIns {
		flag, err := helper.TypeConverter[cdr.FlagSquareItemInfo](item)
		if err != nil {
			log.Println("error converting flag", err)
			return nil, err
		}
		searchFlagInfoReq.FlagId = item.FlagID
		flagDetailReply, err := s.GetFlagDetail(ctx, searchFlagInfoReq)
		if err != nil {
			log.Println("error getting flag detail", err)
			return nil, err
		}
		searchUserReq.UserId = item.UserID
		userInfoReply, err := user.GetClient().GetUserInfo(ctx, searchUserReq)
		if err != nil {
			log.Println("error getting user info", err)
			return nil, err
		}

		flag.UserName = userInfoReply.UserInfo.Name
		flag.SigninId = item.ID
		flag.Content = item.Content
		flag.ChallengeNum = flagDetailReply.Info.ChallengeNum
		// TODO: need to get the number of people who have signed in the siege table
		flag.SiegeNum = 0
		reply.Flags = append(reply.Flags, flag)
	}
	return reply, nil
}

// TODO: need to implement search flag

func startCronTabJob() {
	c := cron.New()
	c.AddFunc("0 16 * * *", UpdateAllFlagStatus)
	c.Start()
}

func UpdateAllFlagStatus() {
	log.Println("start to update all flag status")
	txno := database.GetDB()
	var flags []model.FlagInfo
	err := txno.Model(&model.FlagInfo{}).Where("status in ('running','pending','resurrect')").Find(&flags).Error
	if err != nil {
		log.Println("error getting flags", err)
		return
	}

	for _, flag := range flags {
		txn := txno.Begin()
		var err error
		switch flag.Status {
		case "running":
			err = updateRunningFlagStatus(txn, &flag)
		case "pending":
			err = updatePendingFlagStatus(txn, &flag)
		case "resurrect":
			err = updateResurrectFlagStatus(txn, &flag)
		}
		helper.TransactionHandle(txn, &err)
	}
}

func updateRunningFlagStatus(txn *gorm.DB, flag *model.FlagInfo) error {
	startZeroTime := time.UnixMicro(flag.StartTime).Add(-8 * time.Hour)
	currentSigninNum := len(getSignInlist(txn, flag.ID))
	if startZeroTime.Add(24*time.Hour*time.Duration(currentSigninNum)+1).Sub(time.Now()) >= 0 {
		return nil
	}
	resurrectUsedNum, err := getResurrectUsedNum(txn, flag.ID)
	if err != nil {
		log.Println("error getting resurrect used num", err)
		return err
	}
	if int64(flag.TotalResurrectNum) > resurrectUsedNum {
		log.Println("flag can be resurrected, flagId: ", flag.ID)
		return updateStatusToResurrect(txn, flag)
	}
	flag.Status = "failed"
	err = txn.Model(&model.FlagInfo{}).Where("id = ?", flag.ID).Save(flag).Error
	if err != nil {
		log.Println("error updating flag status", err)
		return err
	}
	return nil
}

func updatePendingFlagStatus(txn *gorm.DB, flag *model.FlagInfo) error {
	startZeroTime := time.UnixMicro(flag.StartTime).Add(-8 * time.Hour)
	nextSignInTine := startZeroTime.Sub(time.Now().Add(24 * time.Hour))
	if nextSignInTine < 0 && nextSignInTine > 24*time.Hour {
		log.Println("no need to update pending flag, flagId: ", flag.ID)
		return nil
	}
	flag.Status = "running"
	err := txn.Model(&model.FlagInfo{}).Where("id = ?", flag.ID).Save(flag).Error
	if err != nil {
		log.Println("error updating flag status", err)
		return err
	}
	return nil
}

func updateResurrectFlagStatus(txn *gorm.DB, flag *model.FlagInfo) error {
	startZeroTime := time.UnixMicro(flag.StartTime).Add(-8 * time.Hour)
	currentSigninNum := len(getSignInlist(txn, flag.ID))
	if startZeroTime.Add(24*time.Hour*time.Duration(currentSigninNum)).Sub(time.Now()) >= 0 {
		flag.Status = "running"
	} else {
		flag.Status = "failed"
	}
	err := txn.Model(&model.FlagInfo{}).Where("id = ?", flag.ID).Save(flag).Error
	if err != nil {
		log.Println("error updating flag status", err)
		return err
	}
	return nil
}
