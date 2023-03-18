package api

import (
	"context"
	"log"

	iampb "github.com/leepala/OldGeneralBackend/Proto/iam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcIAMAddress = "serviceiam.oldgeneral.svc.cluster.local:50051"
	grpcIamClient  iampb.IamClient
)

func getIamClient() iampb.IamClient {
	if grpcIamClient != nil {
		return grpcIamClient
	}
	conn, err := grpc.Dial(grpcIAMAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to dataApi server: %v with error: %v ", grpcIAMAddress, err)
		return nil
	}

	grpcIamClient = iampb.NewIamClient(conn)
	return grpcIamClient
}

func (s *server) IAMLogin(ctx context.Context, in *iampb.IamLoginRequest) (*iampb.IamLoginReply, error) {
	log.Println("received login request", in.RequestId, in.UserName)
	return getIamClient().IAMLogin(ctx, in)
}

func (s *server) IAMRegister(ctx context.Context, in *iampb.CreateUserRequest) (*iampb.CreateUserReply, error) {
	log.Println("received register request", in.RequestId, in.UserName)
	return getIamClient().IAMRegister(ctx, in)
}

func (s *server) IAMCheckLoginStatus(ctx context.Context, in *iampb.IamCheckStatusRequest) (*iampb.IamCheckStatusReply, error) {
	log.Println("received check login status request", in.RequestId)
	return getIamClient().IAMCheckLoginStatus(ctx, in)
}
