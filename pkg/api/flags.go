package api

import (
	"context"
	"log"

	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	"github.com/leepala/OldGeneralBackend/pkg/flags"
)

func (s *server) CreateFlag(ctx context.Context, in *flagspb.CreateFlagRequest) (*flagspb.CreateFlagReply, error) {
	log.Println("create flag request", in.RequestId, in.Info.Id)
	return flags.GetClient().CreateFlag(ctx, in)
}
