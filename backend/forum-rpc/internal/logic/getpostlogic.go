package logic

import (
	"context"
	"errors"
	"time"

	"common/cache"
	"forum-rpc/internal/model"
	"forum-rpc/internal/svc"
	"forum-rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostLogic) GetPost(in *pb.GetPostReq) (*pb.GetPostResp, error) {
	if in.GetId() <= 0 {
		return nil, errors.New("id is required")
	}

	cacheKey := cache.BuildKey("forum:post:detail", in.GetId())
	if cache.IsAvailable() {
		var cached pb.GetPostResp
		if err := cache.GetJSON(l.ctx, cacheKey, &cached); err == nil && cached.Post != nil {
			return &cached, nil
		}
	}

	var post model.ForumPost
	if err := l.svcCtx.DB.Where("id = ? AND status = ?", in.GetId(), "published").First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	_ = l.svcCtx.DB.Model(&model.ForumPost{}).
		Where("id = ?", post.ID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
	post.ViewCount += 1

	resp := &pb.GetPostResp{
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
	}

	if cache.IsAvailable() {
		_ = cache.Set(l.ctx, cacheKey, resp, cache.DefaultTTL)
	}

	return resp, nil
}
