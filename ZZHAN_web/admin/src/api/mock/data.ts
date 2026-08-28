/**
 * 后台 Mock 数据 — 复刻前台数据集形状 + 后台特有（settings/dashboard/profile）。
 */

export const ADMIN_PROFILE = {
  nickname: '阿轩',
  avatar: '',
  email: 'admin@codethink.dev',
}

export const SETTINGS: Record<string, unknown> = {
  blog_name: '小猫的个人博客',
  blog_desc: '记录代码，分享技术，持续成长',
  logo_text: '猫',
  author_name: '阿轩',
  author_role: 'Gopher · 后端工程师',
  author_intro: '一名专注后端开发的全栈工程师，深耕 Go / MySQL / Redis / 云原生，喜欢把复杂的问题讲简单。',
  github: 'https://github.com/yourname',
  email: 'hello@codethink.dev',
  // 关于页 / 首页 hero 联动字段
  tagline: '记录代码，分享技术，持续成长',
  motto: '「写代码是跟计算机对话，写博客是跟自己对话。」',
  location: '北京',
  since: 2019,
  avatar: '',
  // 社交链接
  socials: [
    { name: 'GitHub', icon: 'github', url: 'https://github.com/yourname' },
    { name: '掘金', icon: 'diamond', url: 'https://juejin.cn' },
    { name: 'RSS', icon: 'rss', url: '#' },
  ],
  // 首页终端打字机（每行 `类型|文本`：tk 蓝 / cm 灰 / fn 紫；纯文本默认 tk）
  hero_terminal: [
    'tk|go run server.go',
    'cm|# 今天也在认真写代码',
    'fn|[小猫] ',
    'tk|server listening on :8080',
    'tk|curl http://localhost:8080/api/ping',
    'fn|{"pong":true, ',
    'fn|"uptime":"3d 14h"}',
    'tk|git push origin main ✓ done',
  ].join('\n'),
}

export interface MockArticle {
  id: number
  slug: string
  title: string
  summary: string
  cover_image: string
  category: string
  tags: string[]
  status: 'published' | 'draft' | 'down'
  views: number
  date: string
  updated: string
  content?: string
}

export const MOCK_ARTICLES: MockArticle[] = [
  { id: 1, slug: 'go-concurrency-in-depth', title: 'Go 并发模型深度解析：从 goroutine 到 channel', summary: '并发是 Go 语言最鲜明的名片。本文从 goroutine 的调度模型讲起，深入 channel 与 select。', cover_image: '', category: 'Go', tags: ['Go', '并发', '面试'], status: 'published', views: 12840, date: '2026-08-18', updated: '2026-08-22' },
  { id: 2, slug: 'mysql-index-optimization', title: 'MySQL 索引原理与优化实战：从 B+Tree 到执行计划', summary: '为什么明明建了索引却不生效？手把手读懂 EXPLAIN 执行计划。', cover_image: '', category: 'MySQL', tags: ['MySQL', 'MySQL索引', '性能优化'], status: 'published', views: 15320, date: '2026-08-12', updated: '2026-08-15' },
  { id: 3, slug: 'redis-cache-problems', title: 'Redis 缓存穿透、击穿、雪崩：成因分析与终极解决方案', summary: '缓存三大难题的成因与布隆过滤器、互斥锁、逻辑过期等一揽子方案。', cover_image: '', category: 'Redis', tags: ['Redis', '缓存', '并发'], status: 'published', views: 9860, date: '2026-08-05', updated: '2026-08-09' },
  { id: 4, slug: 'gin-gorm-restful', title: '用 Gin + GORM 搭建一个优雅的 RESTful API 服务', summary: '从项目结构、分层设计到中间件与统一响应。', cover_image: '', category: 'Gin', tags: ['Gin', 'GORM', 'Go'], status: 'published', views: 7420, date: '2026-07-28', updated: '2026-07-30' },
  { id: 5, slug: 'websocket-push-service', title: 'WebSocket 实战：Go 实现实时消息推送服务', summary: '从握手原理到心跳保活，实现广播与断线重连。', cover_image: '', category: 'WebSocket', tags: ['WebSocket', 'Go'], status: 'published', views: 5630, date: '2026-07-15', updated: '2026-07-18' },
  { id: 6, slug: 'sse-streaming-response', title: 'SSE 与流式响应：把 AI 聊天接入你的 Go 服务', summary: '用标准库手写 SSE 服务，AI 流式输出绝配。', cover_image: '', category: 'SSE', tags: ['SSE', 'Go'], status: 'published', views: 4120, date: '2026-07-02', updated: '2026-07-05' },
  { id: 7, slug: 'docker-compose-practice', title: 'Docker 入门到实践：从镜像到 Compose 编排', summary: '镜像、容器、数据卷、网络、Compose 一条龙。', cover_image: '', category: 'Docker', tags: ['Docker', 'Linux'], status: 'published', views: 6840, date: '2026-06-20', updated: '2026-06-22' },
  { id: 8, slug: 'linux-troubleshooting', title: 'Linux 服务排查手册：CPU、内存、磁盘与网络', summary: '一份可直接照着做的线上排查清单。', cover_image: '', category: 'Linux', tags: ['Linux', '性能优化'], status: 'published', views: 5890, date: '2026-06-05', updated: '2026-06-08' },
  { id: 9, slug: 'git-workflow-best-practices', title: 'Git 工作流最佳实践：提交规范、分支策略与 Code Review', summary: '好的 Git 规范能让团队协作事半功倍。', cover_image: '', category: 'Git', tags: ['Git'], status: 'published', views: 3340, date: '2026-05-18', updated: '2026-05-20' },
  { id: 10, slug: 'lru-cache-implementation', title: '手写 LRU 缓存：数据结构与算法实战', summary: '哈希表 + 双向链表实现线程安全 LRU。', cover_image: '', category: '数据结构与算法', tags: ['数据结构', '算法'], status: 'published', views: 4780, date: '2025-05-02', updated: '2025-05-04' },
  { id: 11, slug: 'os-process-thread', title: '操作系统进程与线程：从内核角度看并发', summary: '从内核调度、上下文切换、地址空间讲透并发基础。', cover_image: '', category: '计算机基础', tags: ['计算机基础', '并发'], status: 'published', views: 3960, date: '2025-04-15', updated: '2025-04-16' },
  { id: 12, slug: 'jwt-auth-go', title: 'JWT 认证原理与 Go 实现：从签名算法到刷新令牌', summary: '三段式结构、RS256 与刷新令牌流程。', cover_image: '', category: 'Go', tags: ['JWT', 'JWT认证', 'Go'], status: 'published', views: 6210, date: '2025-03-28', updated: '2025-03-30' },
  { id: 13, slug: 'distributed-transaction', title: '分布式事务入门：从 2PC 到 Seata', summary: '草稿：分布式事务的核心概念梳理。', cover_image: '', category: '计算机基础', tags: ['分布式'], status: 'draft', views: 0, date: '2026-08-20', updated: '2026-08-21' },
  { id: 14, slug: 'go-escape-analysis', title: 'Go 内存模型与逃逸分析', summary: '已下架：栈逃逸与 GC 压力分析。', cover_image: '', category: 'Go', tags: ['Go', '性能优化'], status: 'down', views: 1820, date: '2026-03-10', updated: '2026-04-01' },
]

export const MOCK_CATEGORIES = [
  { id: 1, name: 'Go', slug: 'go', icon: 'code-2', desc: 'Go 语言基础、并发模型与工程实践', color: '#3b82f6', status: 'active', count: 2, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 2, name: 'MySQL', slug: 'mysql', icon: 'database', desc: 'MySQL 原理、索引优化与实战调优', color: '#38bdf8', status: 'active', count: 1, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 3, name: 'Redis', slug: 'redis', icon: 'zap', desc: 'Redis 核心数据结构与缓存设计', color: '#fb7185', status: 'active', count: 1, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 4, name: 'Linux', slug: 'linux', icon: 'terminal', desc: 'Linux 命令、系统管理与排查实战', color: '#a3e635', status: 'active', count: 1, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 5, name: 'Docker', slug: 'docker', icon: 'container', desc: '容器化部署与镜像构建实践', color: '#60a5fa', status: 'active', count: 1, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 6, name: 'Git', slug: 'git', icon: 'git-branch', desc: 'Git 版本控制与团队协作规范', color: '#fb923c', status: 'active', count: 1, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 7, name: 'Gin', slug: 'gin', icon: 'route', desc: 'Gin 框架与 Web 服务开发实战', color: '#34d399', status: 'active', count: 1, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 8, name: 'WebSocket', slug: 'websocket', icon: 'radio', desc: '实时通信协议与长连接实践', color: '#93c5fd', status: 'inactive', count: 1, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 9, name: 'SSE', slug: 'sse', icon: 'activity', desc: '服务端推送与流式响应实现', color: '#fbbf24', status: 'active', count: 1, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 10, name: '数据结构与算法', slug: 'algorithm', icon: 'network', desc: '数据结构、算法设计与刷题笔记', color: '#f472b6', status: 'active', count: 1, created_at: '2024-01-01', updated_at: '2026-08-18' },
  { id: 11, name: '计算机基础', slug: 'cs', icon: 'cpu', desc: '操作系统、网络与计算机组成原理', color: '#94a3b8', status: 'inactive', count: 2, created_at: '2024-01-01', updated_at: '2026-08-18' },
]

export const MOCK_TAGS = [
  { id: 1, name: 'Go', status: 'active', count: 3, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 2, name: '并发', status: 'active', count: 4, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 3, name: 'Gin', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 4, name: 'GORM', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 5, name: 'MySQL索引', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 6, name: 'JWT', status: 'active', count: 2, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 7, name: 'JWT认证', status: 'inactive', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 8, name: 'Redis', status: 'active', count: 2, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 9, name: '缓存', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 10, name: 'WebSocket', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 11, name: 'SSE', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 12, name: 'Docker', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 13, name: 'Linux', status: 'active', count: 2, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 14, name: 'Git', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 15, name: '数据结构', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 16, name: '算法', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 17, name: '分布式', status: 'inactive', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 18, name: '性能优化', status: 'active', count: 3, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 19, name: '面试', status: 'active', count: 2, created_at: '2024-01-05', updated_at: '2024-01-05' },
  { id: 20, name: '源码阅读', status: 'active', count: 1, created_at: '2024-01-05', updated_at: '2024-01-05' },
]

export const MOCK_COMMENTS = [
  { id: 1, article_id: 1, article_title: 'Go 并发模型深度解析：从 goroutine 到 channel', user_name: 'Kevin_z', avatar: '#3b82f6', content: '关于 GMP 模型的讲解很清晰，P 与 M 的绑定关系终于看懂了。', ip: '223.104.***.12', time: '2026-08-22 14:32', status: 'normal' },
  { id: 2, article_id: 2, article_title: 'MySQL 索引原理与优化实战：从 B+Tree 到执行计划', user_name: '小鹿乱撞', avatar: '#fb7185', content: '延迟关联优化分页这个案例太实用了。', ip: '183.14.***.87', time: '2026-08-16 09:15', status: 'normal' },
  { id: 3, article_id: 3, article_title: 'Redis 缓存穿透、击穿、雪崩：成因分析与终极解决方案', user_name: 'Leo_Ma', avatar: '#fbbf24', content: '建议补充布隆过滤器的误判率怎么控制。', ip: '120.36.***.5', time: '2026-08-10 21:47', status: 'normal' },
  { id: 4, article_id: 4, article_title: '用 Gin + GORM 搭建一个优雅的 RESTful API 服务', user_name: '匿名访客', avatar: '#94a3b8', content: '项目分层很规范，能出一期关于 GORM 关联查询的教程吗？', ip: '39.156.***.44', time: '2026-07-29 10:02', status: 'normal' },
  { id: 5, article_id: 5, article_title: 'WebSocket 实战：Go 实现实时消息推送服务', user_name: '阿伟', avatar: '#34d399', content: '心跳保活部分讲得好，之前就踩过连接假死的坑。', ip: '117.136.***.29', time: '2026-07-16 16:38', status: 'normal' },
  { id: 6, article_id: 6, article_title: 'SSE 与流式响应：把 AI 聊天接入你的 Go 服务', user_name: '网络喷子123', avatar: '#ef4444', content: '写得什么垃圾，一点都不如 WebSocket！', ip: '45.207.***.66', time: '2026-07-03 08:20', status: 'banned' },
  { id: 7, article_id: 12, article_title: 'JWT 认证原理与 Go 实现：从签名算法到刷新令牌', user_name: 'Ada_L', avatar: '#c084fc', content: 'RS256 和 HS256 的区别终于弄明白了。', ip: '58.246.***.101', time: '2026-07-12 19:55', status: 'normal' },
  { id: 8, article_id: 7, article_title: 'Docker 入门到实践：从镜像到 Compose 编排', user_name: 'Docker小白', avatar: '#60a5fa', content: '多阶段构建太强了，镜像从 1.2G 降到 18M，已收藏！', ip: '113.87.***.7', time: '2026-06-21 11:30', status: 'normal' },
]

export const DASHBOARD_STATS = [
  { key: 'posts', label: '文章总数', value: 128, trend: '+12', up: true, icon: 'Document', color: '#5b9df6' },
  { key: 'categories', label: '分类数量', value: 11, trend: '+1', up: true, icon: 'Folder', color: '#38bdf8' },
  { key: 'tags', label: '标签数量', value: 46, trend: '+4', up: true, icon: 'PriceTag', color: '#a78bfa' },
  { key: 'comments', label: '评论数量', value: 892, trend: '+36', up: true, icon: 'ChatDotRound', color: '#f59e0b' },
  { key: 'today', label: '今日访问量', value: 4820, trend: '+8.2%', up: true, icon: 'View', color: '#10b981' },
  { key: 'total', label: '总访问量', value: '1,286,420', trend: '+15.6%', up: true, icon: 'TrendCharts', color: '#ef4444' },
]

export const DASHBOARD_CHARTS = {
  week_visits: [
    { d: '08-16', pv: 3210, uv: 2140 }, { d: '08-17', pv: 2860, uv: 1890 },
    { d: '08-18', pv: 4120, uv: 2670 }, { d: '08-19', pv: 3650, uv: 2410 },
    { d: '08-20', pv: 3980, uv: 2550 }, { d: '08-21', pv: 5210, uv: 3280 },
    { d: '08-22', pv: 4820, uv: 3050 },
  ],
  post_trend: [
    { m: '9月', n: 4 }, { m: '10月', n: 6 }, { m: '11月', n: 5 }, { m: '12月', n: 8 },
    { m: '1月', n: 7 }, { m: '2月', n: 4 }, { m: '3月', n: 9 }, { m: '4月', n: 11 },
    { m: '5月', n: 8 }, { m: '6月', n: 12 }, { m: '7月', n: 10 }, { m: '8月', n: 14 },
  ],
  cat_dist: MOCK_CATEGORIES.slice(0, 6).map((c) => ({ name: c.name, value: c.count, color: c.color })),
}

export const DASHBOARD_ARTICLES = {
  recent_posts: MOCK_ARTICLES.slice(0, 5).map((a) => ({ id: a.id, title: a.title, category: a.category, date: a.date, views: a.views })),
  hot_posts: [...MOCK_ARTICLES]
    .sort((a, b) => b.views - a.views)
    .slice(0, 5)
    .map((a) => ({ id: a.id, title: a.title, views: a.views })),
}

export const DASHBOARD_COMMENTS = MOCK_COMMENTS.slice(0, 5).map((c) => ({
  id: c.id,
  user_name: c.user_name,
  avatar: c.avatar,
  content: c.content,
  time: c.time,
}))

export const DASHBOARD_OPERATIONS = [
  { time: '2026-08-22 10:30', user: '阿轩', action: '发布文章', target: '《Go 并发模型深度解析》' },
  { time: '2026-08-22 09:12', user: '小鹿', action: '审核评论', target: '文章 #3 的一条评论' },
  { time: '2026-08-21 18:40', user: '阿轩', action: '更新系统设置', target: 'SEO 描述' },
  { time: '2026-08-21 15:05', user: 'Leo', action: '新建分类', target: '「SSE」' },
  { time: '2026-08-20 22:01', user: 'Ada', action: '删除评论', target: '违规评论 #6' },
  { time: '2026-08-20 10:44', user: '小鹿', action: '保存草稿', target: '《分布式事务入门》' },
  { time: '2026-08-19 09:12', user: 'Leo', action: '禁用用户', target: 'Kevin' },
]
