package flags

import (
	"log"

	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcFlagsAddress = "serviceflags.oldgeneral.svc.cluster.local" + listenPort
	grpcFlagsClient  *flagspb.FlagsClient
)

func GetClient() flagspb.FlagsClient {
	if grpcFlagsClient != nil {
		return *grpcFlagsClient
	}
	conn, err := grpc.Dial(grpcFlagsAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to dataApi server: %v with error: %v ", grpcFlagsAddress, err)
		return nil
	}

	client := flagspb.NewFlagsClient(conn)
	grpcFlagsClient = &client
	return *grpcFlagsClient
}
