package api

import (
	"context"
	"log"
	"net"
	"os"

	apipb "github.com/leepala/OldGeneralBackend/Proto/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	apipb.UnimplementedApiServer
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, exists := metadata.FromIncomingContext(ctx)
	if exists {
		token := md.Get(CONTEXT_USER_TOKEN_AUTHORIZATION_STR)
		if len(token) > 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, CONTEXT_USER_TOKEN_AUTHORIZATION_STR, token[0])
		}
	} else {
		log.Println("cannot get authorization")
	}
	return handler(ctx, req)
}

func StartAndListen() {
	var listenPort = os.Getenv("ListenPort")
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	interceptor := grpc.UnaryInterceptor(unaryInterceptor)
	s := grpc.NewServer(interceptor)
	apipb.RegisterApiServer(s, &server{})
	log.Println("API Server is listening on port 30001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
