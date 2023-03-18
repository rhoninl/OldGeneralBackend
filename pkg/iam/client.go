package iam

import (
	"log"

	iampb "github.com/leepala/OldGeneralBackend/Proto/iam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcIAMAddress = "serviceiam.oldgeneral.svc.cluster.local" + listenPort
	grpcIamClient  *iampb.IamClient
)

func GetClient() iampb.IamClient {
	if grpcIamClient != nil {
		return *grpcIamClient
	}
	conn, err := grpc.Dial(grpcIAMAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to dataApi server: %v with error: %v ", grpcIAMAddress, err)
		return nil
	}

	client := iampb.NewIamClient(conn)
	grpcIamClient = &client
	return *grpcIamClient
}
