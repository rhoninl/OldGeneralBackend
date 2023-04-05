package flags

import (
	"context"
	"log"
	"time"

	"github.com/leepala/OldGeneralBackend/Proto/cdr"
	flagspb "github.com/leepala/OldGeneralBackend/Proto/flags"
	userpb "github.com/leepala/OldGeneralBackend/Proto/user"
	"github.com/leepala/OldGeneralBackend/pkg/database"
	"github.com/leepala/OldGeneralBackend/pkg/helper"
	"github.com/leepala/OldGeneralBackend/pkg/model"
	"github.com/leepala/OldGeneralBackend/pkg/user"
)

func (s *server) PostComment(ctx context.Context, in *flagspb.PostCommentRequest) (*flagspb.PostCommentReply, error) {
	log.Println("post comment request", in)
	comment, err := helper.TypeConverter[model.Comment](in.Comment)
	if err != nil {
		log.Println("error converting comment", err)
		return nil, err
	}

	comment.UserID = in.Comment.UserInfo.Id
	err = database.GetDB().Model(&comment).Save(&comment).Error
	if err != nil {
		log.Println("error saving comment", err)
		return nil, err
	}

	var reply = &flagspb.PostCommentReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
	}
	return reply, nil
}

func (s *server) FetchComment(ctx context.Context, in *flagspb.FetchCommentRequest) (*flagspb.FetchCommentReply, error) {
	log.Println("fetch comment request", in)
	var lastCommentTimeStamp int64 = time.Now().UnixMicro() + 1
	if in.LastCommentId != "" {
		var lastSignInId model.FlagInfo
		err := database.GetDB().Model(&model.Comment{}).Where("id = ?", in.LastCommentId).Find(&lastSignInId).Error
		if err != nil {
			log.Println("error getting last sign in id", err)
			return nil, err
		}
		lastCommentTimeStamp = lastSignInId.CreatedAt
	}
	var comments []model.Comment
	err := database.GetDB().Model(&model.Comment{}).Where("signin_id = ?", in.SigninId).Where("created_at < ?", lastCommentTimeStamp).Order("created_at desc").Limit(int(in.PageSize)).Find(&comments).Error
	if err != nil {
		log.Println("error getting comments", err)
		return nil, err
	}

	searchUserReq := &userpb.GetUserInfoRequest{
		RequestId:   in.RequestId,
		RequestTime: in.RequestTime,
	}
	var commentpbs []*cdr.CommentInfo
	for _, comment := range comments {
		commentpb, err := helper.TypeConverter[cdr.CommentInfo](comment)
		if err != nil {
			log.Println("error converting comment", err)
			return nil, err
		}
		searchUserReq.UserId = comment.UserID
		reply, err := user.GetClient().GetUserInfo(ctx, searchUserReq)
		if err != nil {
			log.Println("error getting user info", err)
			return nil, err
		}
		commentpb.UserInfo = reply.UserInfo
		commentpbs = append(commentpbs, commentpb)
	}

	var reply = &flagspb.FetchCommentReply{
		RequestId: in.RequestId,
		ReplyTime: time.Now().UnixMicro(),
		Comments:  commentpbs,
	}
	return reply, nil
}
