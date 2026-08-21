/**
 * 评论接口：列表 / 发表 / 评论点赞。
 */
import { request } from './http'
import type { CommentDraft, CommentItem, CommentLikeResult, CommentPostResult } from '@/types/models'
import type { Paged } from '@/types/api'

/** 评论列表 GET /articles/{slug}/comments */
export function getComments(slug: string, params?: { page?: number; page_size?: number }): Promise<Paged<CommentItem>> {
  return request<Paged<CommentItem>>({ method: 'GET', url: `/articles/${slug}/comments`, params })
}

/** 发表评论 POST /articles/{slug}/comments */
export function postComment(slug: string, draft: CommentDraft): Promise<CommentPostResult> {
  return request<CommentPostResult>({ method: 'POST', url: `/articles/${slug}/comments`, data: draft })
}

/** 评论点赞 POST /comments/{id}/like（幂等切换） */
export function toggleCommentLike(id: number | string): Promise<CommentLikeResult> {
  return request<CommentLikeResult>({ method: 'POST', url: `/comments/${id}/like` })
}
