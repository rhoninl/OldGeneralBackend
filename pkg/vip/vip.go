package vip

import (
	"context"
	"log"
	"net"
	"time"

	vippb "github.com/leepala/OldGeneralBackend/Proto/vip"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	uuid "github.com/satori/go.uuid"

	"google.golang.org/grpc"
)

const (
	listenPort = ":30001"
)

type server struct {
	vippb.UnimplementedVipServer
}

func StartAndListen() {
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	vippb.RegisterVipServer(s, &server{})
	log.Println("VIP Server is listening on port", listenPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) ChargeVip(ctx context.Context, in *vippb.ChargeVipRequest) (*vippb.ChargeVipReply, error) {
	log.Println("ChargeVip request", in.RequestId, in.UserId)
	var vipInfo = &model.Vip{}
	err := database.GetDB().Model(&vipInfo).Where("user_id = ?", in.UserId).Scan(&vipInfo).Error
	if err != nil {
		log.Printf("cannot get vip info by user id: %s, error: %s", in.UserId, err)
		return nil, err
	}
	if vipInfo.StartTime == 0 {
		vipInfo.ID = uuid.NewV4().String()
		vipInfo.UserID = in.UserId
	}
	if vipInfo.EndTime < time.Now().UnixMicro() {
		vipInfo.StartTime = time.Now().UnixMicro()
		vipInfo.EndTime = vipInfo.StartTime
	}
	vipInfo.EndTime += in.ChargeDuration
	err = database.GetDB().Model(&vipInfo).Where("user_id = ?", in.UserId).Save(&vipInfo).Error
	if err != nil {
		log.Printf("cannot save vip info by user id: %s, error: %s", in.UserId, err)
		return nil, err
	}

	var reply = &vippb.ChargeVipReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}

func (s *server) GetVipStatus(ctx context.Context, in *vippb.GetVipStatusRequest) (*vippb.GetVipStatusReply, error) {
	log.Println("GetVipStatus request", in.RequestId, in.UserId)
	var vipInfo = &model.Vip{}
	err := database.GetDB().Model(&vipInfo).Where("user_id = ?", in.UserId).Find(&vipInfo).Error
	if err != nil {
		log.Printf("cannot get vip info by user id: %s, error: %s", in.UserId, err.Error())
		return nil, err
	}
	var reply = &vippb.GetVipStatusReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		EndTime:   vipInfo.EndTime,
	}
	return reply, nil
}
