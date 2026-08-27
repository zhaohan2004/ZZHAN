/**
 * 评论接口：列表 / 发表 / 回复列表 / 评论点赞。
 */
import { request } from './http'
import type { CommentDraft, CommentItem, CommentLikeResult, CommentPostResult } from '@/types/models'
import type { Paged } from '@/types/api'

/** 评论列表 GET /comments/{slug} */
export function getComments(slug: string, params?: { page?: number; page_size?: number }): Promise<Paged<CommentItem>> {
  return request<Paged<CommentItem>>({ method: 'GET', url: `/comments/${slug}`, params })
}

/** 回复列表 GET /comments/replies/{id} */
export function getReplies(id: number | string, params?: { page?: number; page_size?: number }): Promise<Paged<CommentItem>> {
  return request<Paged<CommentItem>>({ method: 'GET', url: `/comments/replies/${id}`, params })
}

/** 发表评论 POST /comments/{slug} */
export function postComment(slug: string, draft: CommentDraft): Promise<CommentPostResult> {
  return request<CommentPostResult>({ method: 'POST', url: `/comments/${slug}`, data: draft })
}

/** 评论点赞 POST /like/comment/{id}（幂等切换） */
export function toggleCommentLike(id: number | string): Promise<CommentLikeResult> {
  return request<CommentLikeResult>({ method: 'POST', url: `/like/comment/${id}` })
}
