package user

import (
	"log"

	userpb "github.com/leepala/OldGeneralBackend/Proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcUserAddress = "serviceuser.oldgeneral.svc.cluster.local" + listenPort
)

var userClient *userpb.UserClient

func GetClient() userpb.UserClient {
	if userClient != nil {
		return *userClient
	}
	conn, err := grpc.Dial(grpcUserAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to User server: %v with error: %v ", grpcUserAddress, err)
		return nil
	}

	client := userpb.NewUserClient(conn)
	userClient = &client
	return *userClient
}
