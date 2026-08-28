package dto

// ArticleLikeResult 文章点赞响应
type ArticleLikeResult struct {
	Liked bool  `json:"liked"`
	Likes int32 `json:"likes"`
}

// CommentLikeResult 评论点赞响应
type CommentLikeResult struct {
	Liked     bool  `json:"liked"`
	LikeCount int32 `json:"like_count"`
}
