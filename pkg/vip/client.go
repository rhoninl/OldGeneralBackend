package vip

import (
	"log"

	vippb "github.com/leepala/OldGeneralBackend/Proto/vip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcvipAddress = "servicevip.oldgeneral.svc.cluster.local" + listenPort
)

var vipClient *vippb.VipClient

func GetClient() vippb.VipClient {
	if vipClient != nil {
		return *vipClient
	}
	conn, err := grpc.Dial(grpcvipAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to vip server: %v with error: %v ", grpcvipAddress, err)
		return nil
	}

	client := vippb.NewVipClient(conn)
	vipClient = &client
	return *vipClient
}
