package logic

import (
	"context"
	"errors"
	"strings"
	"time"

	"forum-rpc/internal/model"
	"forum-rpc/internal/svc"
	"forum-rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePostLogic) CreatePost(in *pb.CreatePostReq) (*pb.CreatePostResp, error) {
	if in.GetBoardId() <= 0 {
		return nil, errors.New("board_id is required")
	}
	if in.GetUserId() <= 0 {
		return nil, errors.New("user_id is required")
	}
	title := strings.TrimSpace(in.GetTitle())
	content := strings.TrimSpace(in.GetContent())
	if title == "" {
		return nil, errors.New("title is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}

	post := model.ForumPost{
		BoardID:  uint(in.GetBoardId()),
		UserID:   uint(in.GetUserId()),
		Title:    title,
		Content:  content,
		Status:   "published",
		IsPinned: false,
		IsLocked: false,
	}
	if err := l.svcCtx.DB.Create(&post).Error; err != nil {
		return nil, err
	}

	return &pb.CreatePostResp{
		Post: &pb.PostItem{
			Id:           int64(post.ID),
			BoardId:      int64(post.BoardID),
			UserId:       int64(post.UserID),
			Title:        post.Title,
			Content:      post.Content,
			Status:       post.Status,
			IsPinned:     post.IsPinned,
			IsLocked:     post.IsLocked,
			ViewCount:    post.ViewCount,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			CreatedAt:    post.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}
