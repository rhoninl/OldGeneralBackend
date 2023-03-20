package flags

import (
	"context"
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
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

const (
	listenPort = ":30001"
)

type server struct {
	flagspb.UnimplementedFlagsServer
}

func StartAndListen() {
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	flagspb.RegisterFlagsServer(s, &server{})
	log.Println("API Server is listening on port", listenPort)
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
	flag.ID = uuid.NewV4().String()
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
		reply.Flags = append(reply.Flags, f)
	}
	return reply, nil
}

func (s *server) GetFlagDetail(ctx context.Context, in *flagspb.GetFlagDetailRequest) (*flagspb.GetFlagDetailReply, error) {
	log.Println("get flag detail request", in.RequestId, in.FlagId)
	var flag model.FlagInfo
	err := database.GetDB().Model(&flag).Where("id = ?", in.FlagId).Find(&flag).Error
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

	f.UserAvatar = userInfoReply.UserInfo.Avatar
	f.UserName = userInfoReply.UserInfo.Name

	var reply = &flagspb.GetFlagDetailReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		Info:      f,
	}
	return reply, nil
}
