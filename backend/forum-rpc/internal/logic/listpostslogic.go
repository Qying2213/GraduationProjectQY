package logic

import (
	"context"
	"time"

	"forum-rpc/internal/model"
	"forum-rpc/internal/svc"
	"forum-rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPostsLogic) ListPosts(in *pb.ListPostsReq) (*pb.ListPostsResp, error) {
	page := in.GetPage()
	pageSize := in.GetPageSize()
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := l.svcCtx.DB.Model(&model.ForumPost{}).Where("status = ?", "published")
	if in.GetBoardId() > 0 {
		query = query.Where("board_id = ?", in.GetBoardId())
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var posts []model.ForumPost
	if err := query.
		Order("is_pinned DESC").
		Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	items := make([]*pb.PostItem, 0, len(posts))
	for _, p := range posts {
		items = append(items, &pb.PostItem{
			Id:           int64(p.ID),
			BoardId:      int64(p.BoardID),
			UserId:       int64(p.UserID),
			Title:        p.Title,
			Content:      p.Content,
			Status:       p.Status,
			IsPinned:     p.IsPinned,
			IsLocked:     p.IsLocked,
			ViewCount:    p.ViewCount,
			LikeCount:    p.LikeCount,
			CommentCount: p.CommentCount,
			CreatedAt:    p.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListPostsResp{
		List:  items,
		Total: total,
	}, nil
}
