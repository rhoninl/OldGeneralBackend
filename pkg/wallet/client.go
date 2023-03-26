package wallet

import (
	"log"

	walletpb "github.com/leepala/OldGeneralBackend/Proto/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcWalletAddress = "servicewallet.oldgeneral.svc.cluster.local" + listenPort
)

var walletClient *walletpb.WalletClient

func GetClient() walletpb.WalletClient {
	if walletClient != nil {
		return *walletClient
	}
	conn, err := grpc.Dial(grpcWalletAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to User server: %v with error: %v ", grpcWalletAddress, err)
		return nil
	}

	client := walletpb.NewWalletClient(conn)
	walletClient = &client
	return *walletClient
}
