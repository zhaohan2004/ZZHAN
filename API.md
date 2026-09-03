# ZZHAN 博客系统 · API 接口文档

> 版本：v1.0（依据后端源码 `ZZHAN/` 逐接口核对整理）
> 技术栈：Go + Gin + GORM + MySQL + Redis；认证：JWT（access_token + refresh_token）
> 适用对象：前台 / 后台前端开发人员
> 说明：文档仅描述源码中真实存在的接口、参数、字段与业务逻辑；源码中无法确定的信息均标注为「源码未明确」。凡标注 **[无]** 表示源码未提供该数据。

---

## 目录

1. [接口总览](#一接口总览)
2. [统一响应格式](#二统一响应格式)
3. [认证机制](#三认证机制)
4. [权限说明](#四权限说明)
5. [错误码说明](#五错误码说明)
6. [业务接口详情](#六业务接口详情)
7. [分页说明](#七分页说明)

---

## 一、接口总览

统一前缀：`/api/v1`。后台路由以 `/admin` 开头；其余为前台路由。

### 前台公开接口

| # | 方法 | 路径 | 认证 | 说明 |
|---|------|------|------|------|
| 1 | GET | `/health` | 无 | 健康检查 |
| 2 | GET | `/site` | 无 | 站点信息 |
| 3 | GET | `/articles` | 无 | 文章列表 |
| 4 | GET | `/articles/:slug` | OptionalAuth | 文章详情 |
| 5 | GET | `/categories` | 无 | 分类列表 |
| 6 | GET | `/tags` | 无 | 标签列表 |
| 7 | GET | `/archives` | 无 | 文章归档 |
| 8 | GET | `/stats` | 无 | 站点统计 |
| 9 | GET | `/comments/:slug` | 无 | 某篇文章的评论列表 |
| 10 | GET | `/comments/replies/:id` | 无 | 某条评论的回复列表 |

### 前台登录用户接口（user）

| # | 方法 | 路径 | 认证 | 说明 |
|---|------|------|------|------|
| 11 | POST | `/auth/github` | 无 | GitHub OAuth 登录 |
| 12 | POST | `/auth/refresh` | 无 | 刷新 access_token |
| 13 | GET | `/auth/me` | user | 当前登录用户信息 |
| 14 | POST | `/auth/logout` | user | 退出登录 |
| 15 | POST | `/comments/:slug` | user | 发表评论 |
| 16 | POST | `/like/article/:slug` | user | 文章点赞 / 取消点赞 |
| 17 | POST | `/like/comment/:id` | user | 评论点赞 / 取消点赞 |

### 后台认证接口（部分公开 / 部分 admin）

| # | 方法 | 路径 | 认证 | 说明 |
|---|------|------|------|------|
| 18 | GET | `/admin/auth/captcha` | 无 | 图形验证码 |
| 19 | POST | `/admin/auth/login` | 无 | 后台管理员登录 |
| 20 | POST | `/admin/auth/refresh` | 无 | 后台刷新 access_token |
| 21 | POST | `/admin/auth/logout` | admin | 后台退出登录 |
| 22 | GET | `/admin/profile` | admin | 管理员资料 |
| 23 | PUT | `/admin/profile` | admin | 更新管理员资料 |

### 后台管理接口（admin）

| # | 方法 | 路径 | 说明 |
|---|------|------|------|
| 24 | GET | `/admin/articles` | 文章列表 |
| 25 | GET | `/admin/articles/:id` | 文章详情 |
| 26 | POST | `/admin/articles` | 新建文章 |
| 27 | PUT | `/admin/articles/:id` | 更新文章 |
| 28 | DELETE | `/admin/articles/:id` | 删除文章 |
| 29 | PUT | `/admin/articles/:id/status` | 修改文章状态 |
| 30 | GET | `/admin/categories` | 分类列表 |
| 31 | GET | `/admin/categories/:id` | 分类详情 |
| 32 | POST | `/admin/categories` | 新建分类 |
| 33 | PUT | `/admin/categories/:id` | 更新分类 |
| 34 | PUT | `/admin/categories/:id/status` | 修改分类状态 |
| 35 | DELETE | `/admin/categories/:id` | 删除分类 |
| 36 | GET | `/admin/tags` | 标签列表 |
| 37 | GET | `/admin/tags/:id` | 标签详情 |
| 38 | POST | `/admin/tags` | 新建标签 |
| 39 | PUT | `/admin/tags/:id` | 更新标签 |
| 40 | PUT | `/admin/tags/:id/status` | 修改标签状态 |
| 41 | DELETE | `/admin/tags/:id` | 删除标签 |
| 42 | GET | `/admin/comments` | 评论列表 |
| 43 | GET | `/admin/comments/:id` | 评论详情 |
| 44 | PUT | `/admin/comments/:id/status` | 修改评论状态 |
| 45 | DELETE | `/admin/comments/:id` | 删除评论 |
| 46 | GET | `/admin/users` | 用户列表 |
| 47 | GET | `/admin/users/:id` | 用户详情 |
| 48 | PUT | `/admin/users/:id/status` | 修改用户状态 |
| 49 | DELETE | `/admin/users/:id` | 删除用户 |
| 50 | GET | `/admin/settings` | 读取系统设置 |
| 51 | PUT | `/admin/settings` | 保存系统设置 |
| 52 | POST | `/upload/image` | 上传图片 |
| 53 | GET | `/admin/dashboard/stats` | 仪表盘统计 |
| 54 | GET | `/admin/dashboard/articles` | 仪表盘-最新文章 |
| 55 | GET | `/admin/dashboard/comments` | 仪表盘-最新评论 |
| 56 | GET | `/admin/dashboard/operations` | 仪表盘-最新操作记录 |
| 57 | GET | `/admin/operation-logs` | 操作日志列表 |

---

## 二、统一响应格式

成功与失败均返回 JSON，字段固定为 `code`、`message`、`data`（成功且无数据时 `data` 可省略或缺省为 `null`）：

```json
{
  "code": 0,
  "message": "成功",
  "data": { }
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 业务码：`0` 表示成功，非 `0` 表示各类错误（见 [错误码说明](#五错误码说明)） |
| message | string | 提示信息，可直接展示给用户 |
| data | object/array/null | 业务数据；接口无数据时该字段不返回或为 `null` |

> **注意**：源码中存在两套 code 风格（详见错误码章节）：多数后台资源管理接口（articles/categories/tags/comments/users/settings 等）经 `pkg/response` 输出，`code` 直接使用 HTTP 状态码（如 400/401/403/404/500）或业务码；而认证类接口（admin auth、web auth）与少数 controller 使用手写 `gin.H`，`code` 形如 `40000 / 40100 / 50000 / 3004`。两套 `code` 均表示「非 0 即失败」，判断成功请统一使用 `code === 0`。

### HTTP 状态码

- 全部成功响应使用 HTTP `200`（无论 `code` 是否非 0）。
- 各接口失败时的具体 HTTP 状态码见「错误情况」小节。

---

## 三、认证机制

### 3.1 认证方式总览

| 认证值 | 含义 | 判定规则 |
|--------|------|----------|
| **无** | 无需任何令牌即可访问 | 直接放行 |
| **OptionalAuth** | 可选认证 | 若请求头带 `Authorization: Bearer <access_token>` 且令牌合法，则将用户信息注入上下文（用于判断点赞等个人状态）；无令牌或令牌非法均不拦截，继续放行 |
| **user** | 前台登录用户 | 必须携带合法 `access_token`，且为 user 类型活跃令牌（单设备登录校验），否则返回 401 |
| **admin** | 后台管理员 | 必须携带合法 `access_token`，且为 admin 类型活跃令牌（单设备登录校验），否则返回 401 |

所有需认证请求统一在请求头携带：

```
Authorization: Bearer <access_token>
```

### 3.2 JWT（JSON Web Token）

- 算法：HS256。
- 签发者（iss）：`ZZHAN`。
- Token 分两种，通过 claim 字段 `token_type` 区分：
  - `access_token`：访问令牌，请求时放入 `Authorization` 头。
  - `refresh_token`：刷新令牌，仅用于换取新的 `access_token`。
- Claims 结构（JWT 载荷）：

```json
{
  "user_id": 1,
  "username": "admin",
  "token_type": "access_token",
  "iss": "ZZHAN",
  "iat": 1700000000,
  "nbf": 1700000000,
  "exp": 1700007200
}
```

| Claim | 类型 | 说明 |
|-------|------|------|
| user_id | int | 用户 / 管理员 ID |
| username | string | 用户名（前台为用户昵称、后台为管理员用户名） |
| token_type | string | `access_token` 或 `refresh_token` |
| iss | string | 固定为 `ZZHAN` |
| iat / nbf / exp | int | 签发时间 / 生效时间 / 过期时间（Unix 秒） |

### 3.3 Access Token 与 Refresh Token

- 登录成功后同时返回 `access_token` 与 `refresh_token`。
- `access_token` 用于访问受保护接口；过期或失效后，客户端使用 `refresh_token` 调用刷新接口获取新的 `access_token`，无需重新登录。
- 有效期由后端配置决定：源码默认配置 `access_expire_hours: 2`（2 小时）、`refresh_expire_hours: 168`（7 天），**具体以服务端实际配置为准（源码未明确是否为固定值）**。
- 刷新成功仅返回新的 `access_token`，不返回新的 `refresh_token`（refresh_token 继续沿用）。

### 3.4 GitHub OAuth 登录（前台）

登录流程（前端负责第 1 步）：

1. 前端跳转 GitHub 授权，用户同意后 GitHub 回调，前端拿到授权临时 `code`。
2. 前端调用 `POST /api/v1/auth/github`，携带 `code`（及可选的 `redirect_uri`）。
3. 后端用 `code` + GitHub ClientID/Secret 换取 GitHub `access_token`，再调用 GitHub API 获取用户信息。
4. 后端按 GitHub 用户 ID 查找/创建 `users` 记录（`provider = "github"`，`openid` 为 GitHub 用户 ID 字符串）；已存在用户则同步更新昵称/头像。
5. 签发 JWT token 对返回前端；并更新用户最后登录时间。

> 服务端请求 GitHub API 时若配置了 `github.proxy` 则走该 HTTP 代理（配置项说明，非业务约束）。

### 3.5 单设备登录（Redis 活跃 Token）

- 同一用户（user 或 admin）同一时间仅允许一台设备在线。
- 登录 / 刷新成功后，后端将最新 `access_token`、`refresh_token` 写入 Redis 记为「活跃 token」。
- 请求受保护接口时，`Auth` 中间件会校验当前令牌是否等于该用户最新的活跃令牌；不一致则判定「账号已在其他设备登录」，返回 401。
- 新设备登录时，旧设备的活跃 access_token 与 refresh_token 会先被加入黑名单再被替换，旧设备随即失效。
- Redis 不可用时中间件采取**降级放行**策略（源码注释：其他 Redis 错误按服务故障放行），仅黑名单查询与活跃 token 查询逻辑受此影响。

### 3.6 Redis 黑名单（退出 / 失效）

- 退出登录时，后端把当前 `access_token` 按「剩余有效时长」加入 Redis 黑名单，同时把该用户活跃的 refresh_token 加入黑名单，并清除活跃 token 记录。
- `Auth` 中间件会先查黑名单：命中的令牌直接拒绝（`token 已失效`）。
- 刷新接口同样校验 refresh_token 是否已在黑名单。

### 3.7 后台图形验证码（登录前）

- `GET /admin/auth/captcha` 返回 `captcha_id` 与 base64 图片。
- 登录时必须回传 `captcha_id` 与用户输入的 `captcha`，服务端校验，错误返回「验证码错误 / 图形验证码错误」。

---

## 四、权限说明

系统无 RBAC 角色表，仅两类身份，由 `Auth` 中间件的 `userType` 参数区分：

| 身份 | userType | 覆盖接口 | 获得身份方式 |
|------|----------|----------|--------------|
| 前台登录用户 | `user` | `/auth/me`、`/auth/logout`、发表评论、点赞 | GitHub OAuth 登录 |
| 后台管理员 | `admin` | 全部 `/admin/*`（除公开的 captcha/login/refresh）、`/profile`、`/upload/image` | 账号密码 + 验证码登录（管理员表为单人设计） |
| 匿名（公开） | — | 其余所有读取接口 | — |

- 后台写操作接口额外套用 `OperationLog` 中间件：POST/PUT/DELETE 成功（HTTP 2xx）后自动写入 `operation_logs` 表，记录动作与操作对象（详见「操作日志」模块）。
- `GET /articles/:slug` 使用 `OptionalAuth`：携带合法 user 令牌时返回的 `liked` 字段反映当前用户点赞状态；未登录则 `liked` 恒为 `false`。

---

## 五、错误码说明

> 源码统一业务码定义位于 `pkg/errors`，认证类 controller 另有手写五位数风格码。下表两者均列出。

### 5.1 通用码（`pkg/response` / `pkg/errors` 使用，与 HTTP 状态码一致）

| code | HTTP | 含义 | 典型触发 |
|------|------|------|----------|
| 0 | 200 | 成功 | — |
| 400 | 400 | 请求参数错误 / 业务校验失败 | 绑定参数失败、重复名称、参数非法 |
| 401 | 401 | 未授权 | 缺少令牌、令牌格式错误、黑名单命中、非活跃令牌、未登录 |
| 403 | 403 | 禁止访问 | 无权限 |
| 404 | 404 | 资源不存在 | 按 ID/slug 查无记录 |
| 500 | 500 | 服务器内部错误 | 数据库或内部异常 |

### 5.2 业务码（`pkg/errors` 定义）

| code | 含义 |
|------|------|
| 1001 | 用户不存在 |
| 1002 | 用户已存在 |
| 1003 | 用户名或密码错误 |
| 1004 | 用户已被禁用 |
| 1005 | 无效的令牌 |
| 1006 | 令牌已过期 |
| 2001 | 参数错误 |
| 2002 | 缺少必要参数 |
| 2003 | 参数格式错误 |
| 3001 | 资源不存在 |
| 3002 | 资源已存在 |
| 3003 | 资源已被锁定 |
| 3004 | 图形验证码错误 |

> 说明：`pkg/errors` 定义了上述业务码与对应默认文案，但**各 controller 实际是否使用这些业务码并不一致**；多数直接使用 `pkg/response` 的 `code = HTTP 状态码`，并自定义 `message`。因此前端应以 `code === 0` 判成功，非 0 时优先读取并展示 `message`，无需对具体业务码做分支。

### 5.3 认证类接口的手写码（admin/web auth、site controller）

| code | HTTP | 含义 |
|------|------|------|
| 0 | 200 | 成功 |
| 40000 | 400 | 请求参数错误（如 JSON 绑定失败、缺少 access_token） |
| 40001 | 400 | 请求参数错误（`PUT /admin/settings` 绑定失败时使用） |
| 40100 | 401 | 未登录 / 刷新失败 / 令牌无效 |
| 50000 | 500 | 服务器内部错误（生成验证码失败、登录失败、退出失败、获取/保存失败等） |
| 3004 | 400 | 图形验证码错误（后台登录校验验证码失败时返回，HTTP 400） |

### 5.4 各接口错误情况

下文「业务接口详情」中，每个接口的「错误情况」仅列出**该接口源码中实际出现的分支**；未出现的通用错误不再赘述。

---

## 六、业务接口详情

> 模块顺序：健康检查 → 站点信息 → 前台认证 → 前台文章 → 分类/标签/归档 → 站点统计 → 评论 → 点赞 → 后台认证 → 后台文章管理 → 后台分类管理 → 后台标签管理 → 后台评论管理 → 后台用户管理 → 系统设置 → 文件上传 → 仪表盘 → 操作日志。

---

### 健康检查

#### 1. 健康检查

请求方式：`GET /api/v1/health`

认证方式：无

接口说明：探测服务是否存活。**该接口响应结构与其他接口不同，为固定的 `status/message` 字段，没有 `code`。**

请求参数：无

请求示例：无

成功响应：

```json
{
  "status": "ok",
  "message": "ZZHAN API is running"
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| status | string | 固定为 `ok` |
| message | string | 固定为 `ZZHAN API is running` |

错误情况：无特殊错误分支（服务不可达即网络错误）。

---

### 站点信息（前台公开）

#### 2. 获取站点信息

请求方式：`GET /api/v1/site`

认证方式：无

接口说明：返回博客前台展示所需的站点配置信息（首页 / 关于页等），数据来自 `site_settings` 键值对表。**该接口由 controller 手写 `gin.H` 输出（含 `code/message/data`）。**

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "name": "",
    "tagline": "",
    "bio": "",
    "github": "",
    "email": "",
    "socials": [],
    "author": "",
    "role": "",
    "motto": "",
    "location": "",
    "since": 0,
    "avatar": "",
    "hero_terminal": ""
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| name | string | 站点名（配置键 `blog_name`） |
| tagline | string | 标语（优先取配置键 `tagline`，其次 `blog_desc`） |
| bio | string | 作者简介（配置键 `author_intro`） |
| github | string | GitHub 地址（配置键 `github`） |
| email | string | 邮箱（配置键 `email`） |
| socials | array | 社交链接列表（配置键 `socials` 的 JSON 解析结果） |
| socials[].name | string | 社交名称（源码未明确更多字段含义） |
| socials[].icon | string | 图标名 |
| socials[].url | string | 链接地址 |
| author | string | 作者名（配置键 `author_name`） |
| role | string | 作者角色（配置键 `author_role`） |
| motto | string | 座右铭（配置键 `motto`） |
| location | string | 所在地（配置键 `location`） |
| since | int | 建站年份（配置键 `since`，字符串转整数） |
| avatar | string | 头像 URL（配置键 `avatar`） |
| hero_terminal | string | Hero 终端内容（配置键 `hero_terminal`） |

> 注：`socials` 字段由后端将配置 JSON 解析为对象数组，具体对象字段与 `SiteResponse` 中一致。所有配置均可能为空字符串，取决于后台「系统设置」是否已保存对应键。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 50000 | 500 | 获取站点信息失败 |

---

### 前台认证

#### 3. GitHub OAuth 登录

请求方式：`POST /api/v1/auth/github`

认证方式：无

接口说明：使用 GitHub OAuth 授权码换取登录态。成功返回 `access_token`、`refresh_token` 与用户信息；新用户自动注册。

请求参数（Body，JSON）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| code | string | 是 | body | GitHub 授权临时 code |
| redirect_uri | string | 否 | body | GitHub 回调地址（服务端调用 GitHub 换取 token 时附加） |

请求示例：

```json
{
  "code": "gho_xxxxxxxx",
  "redirect_uri": "https://example.com/auth/github/callback"
}
```

成功响应：

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "access_token": "<jwt>",
    "refresh_token": "<jwt>",
    "user": {
      "id": 1,
      "provider": "github",
      "nickname": "Kevin_z",
      "avatar": "https://avatars.githubusercontent.com/..."
    }
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| access_token | string | 访问令牌，请求受保护接口时放 `Authorization` 头 |
| refresh_token | string | 刷新令牌，用于换取新的 access_token |
| user | object | 用户信息 |
| user.id | int | 用户 ID |
| user.provider | string | 登录方式，前台固定 `github` |
| user.nickname | string | 昵称（GitHub Name，缺省回退 Login） |
| user.avatar | string | 头像 URL |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40000 | 400 | 请求参数错误（JSON 绑定失败） |
| 50000 | 500 | 登录失败：获取 GitHub access_token 失败 / 获取 GitHub 用户信息失败 / 查找或创建用户失败 / 生成 token 失败（message 前缀均为「登录失败：」） |

---

#### 4. 前台刷新令牌

请求方式：`POST /api/v1/auth/refresh`

认证方式：无

接口说明：使用 `refresh_token` 换取新的 `access_token`。校验 refresh_token 合法性、黑名单状态，成功后更新活跃令牌。

请求参数（Body，JSON）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| refresh_token | string | 是 | body | 登录时下发的 refresh_token |

请求示例：

```json
{
  "refresh_token": "<jwt>"
}
```

成功响应：

```json
{
  "code": 0,
  "message": "刷新成功",
  "data": {
    "access_token": "<jwt>"
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| access_token | string | 新的访问令牌 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40000 | 400 | 请求参数错误（JSON 绑定失败） |
| 40100 | 401 | 刷新失败：refresh_token 无效 / 已失效 / 检查黑名单失败 / 生成新 access_token 失败（message 前缀「刷新失败：」） |

---

#### 5. 获取当前登录用户信息

请求方式：`GET /api/v1/auth/me`

认证方式：user

接口说明：返回当前登录用户的信息。

请求参数：无（认证通过后从令牌上下文取用户 ID）

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "id": 1,
    "provider": "github",
    "nickname": "Kevin_z",
    "avatar": "https://avatars.githubusercontent.com/..."
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 用户 ID |
| provider | string | 登录方式（github） |
| nickname | string | 昵称 |
| avatar | string | 头像 URL |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40100 | 401 | 未登录（未取到上下文 user_id） |
| 50000 | 500 | 获取用户信息失败：用户不存在等（message 前缀「获取用户信息失败：」） |

> 注：user 认证失败时，`Auth` 中间件也可能直接返回 401（未携带令牌 / 格式错误 / 黑名单 / 他端登录等），见 [错误码说明](#五错误码说明) 与 [认证机制](#三认证机制)。

---

#### 6. 前台退出登录

请求方式：`POST /api/v1/auth/logout`

认证方式：user

接口说明：将当前 access_token 加入黑名单、清除该用户活跃令牌并使 refresh_token 失效。令牌本身无效或已过期时仍返回成功（幂等）。

请求参数：无（从 `Authorization` 头解析令牌）

> 认证中间件已要求携带合法令牌；controller 直接从 `Authorization` 头去除 `Bearer ` 前缀后取用。

成功响应：

```json
{
  "code": 0,
  "message": "退出成功"
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40000 | 400 | 缺少 access_token（Authorization 头为空） |
| 50000 | 500 | 退出失败（加入黑名单出错，message 前缀「退出失败：」） |

---

### 前台文章

#### 7. 前台文章列表

请求方式：`GET /api/v1/articles`

认证方式：无

接口说明：分页获取已发布（`status = published`）的文章列表，可按分类、标签、关键词筛选。列表项不含正文。

请求参数（Query）：**分页参数必填**

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| page | int | 是 | query | 页码，从 1 开始（`min=1`） |
| size | int | 是 | query | 每页大小（`min=1, max=100`） |
| category_id | int | 否 | query | 按分类 ID 筛选 |
| tag_id | int | 否 | query | 按标签 ID 筛选 |
| keyword | string | 否 | query | 标题/摘要关键词模糊搜索 |

> 注：本接口分页参数名为 `page` 与 `size`（区别于后台的 `page_size`），由源码 `response.PageRequest`（`page` / `size`）决定，且二者均带 `required` 校验，缺省即报「参数错误」。

请求示例：

```
GET /api/v1/articles?page=1&size=10&category_id=1&keyword=golang
```

成功响应（分页结构）：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Go 并发模型",
        "slug": "go-concurrency",
        "summary": "摘要…",
        "cover_image": "https://...",
        "category_id": 1,
        "category_name": "Go",
        "author_name": "阿轩",
        "tags": [ { "id": 1, "name": "Go", "slug": "go" } ],
        "views": 12840,
        "likes": 862,
        "comment_count": 56,
        "published_at": "2026-08-18T00:00:00+08:00"
      }
    ],
    "total": 128,
    "page": 1,
    "size": 10,
    "total_page": 13
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| list | array | 文章列表项（分页外层结构见 [分页说明](#七分页说明)） |
| list[].id | int | 文章 ID |
| list[].title | string | 标题 |
| list[].slug | string | URL 别名 |
| list[].summary | string | 摘要 |
| list[].cover_image | string | 封面图 URL（可为空字符串） |
| list[].category_id | int | 分类 ID |
| list[].category_name | string | 分类名称（关联查询） |
| list[].author_name | string | 作者名称 |
| list[].tags | array | 标签列表 |
| list[].tags[].id | int | 标签 ID |
| list[].tags[].name | string | 标签名称 |
| list[].tags[].slug | string | 标签 URL 别名 |
| list[].views | int | 浏览量 |
| list[].likes | int | 点赞数 |
| list[].comment_count | int | 评论数 |
| list[].published_at | string | 发布时间（RFC3339 时间格式） |
| total | int | 总记录数 |
| page / size / total_page | int | 当前页 / 每页大小 / 总页数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 参数错误：分页或查询参数绑定失败（message 含具体校验错误） |
| 500 | 500 | 获取文章列表失败 |

---

#### 8. 前台文章详情

请求方式：`GET /api/v1/articles/:slug`

认证方式：OptionalAuth

接口说明：按 slug 获取已发布文章详情，含完整 Markdown 正文。携带合法 user 令牌时返回 `liked` 反映当前用户点赞状态；每次成功访问会做浏览量计数（同一 IP 对同一文章去重）。**访问前台可见文章详情由请求方（如浏览器）记录客户端 IP 无需传参。**

请求参数（Path）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| slug | string | 是 | path | 文章 URL 别名 |

请求示例：

```
GET /api/v1/articles/go-concurrency
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "id": 1,
    "title": "Go 并发模型",
    "slug": "go-concurrency",
    "summary": "摘要…",
    "cover_image": "https://...",
    "category_id": 1,
    "category_name": "Go",
    "author_name": "阿轩",
    "tags": [ { "id": 1, "name": "Go", "slug": "go" } ],
    "content": "# Markdown 正文…",
    "views": 12841,
    "likes": 862,
    "liked": false,
    "comment_count": 56,
    "published_at": "2026-08-18T00:00:00+08:00"
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 文章 ID |
| title | string | 标题 |
| slug | string | URL 别名 |
| summary | string | 摘要 |
| cover_image | string | 封面图 URL |
| category_id | int | 分类 ID |
| category_name | string | 分类名称 |
| author_name | string | 作者名称 |
| tags | array | 标签列表（结构同列表接口） |
| content | string | Markdown 正文 |
| views | int | 浏览量 |
| likes | int | 点赞数 |
| liked | bool | 当前用户是否已点赞（OptionalAuth 未登录时为 false） |
| comment_count | int | 评论数 |
| published_at | string | 发布时间 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 参数错误：slug 绑定失败 |
| 404 | 404 | 文章不存在（未查到 status=published 的文章） |
| 500 | 500 | 获取文章详情失败（非「记录不存在」的其他错误） |

---

### 分类 / 标签 / 归档

#### 9. 前台分类列表

请求方式：`GET /api/v1/categories`

认证方式：无

接口说明：返回所有启用（`status = active`）的分类及每个分类下的已发布文章数，按排序值升序。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": [
    {
      "id": 1,
      "name": "Go",
      "slug": "go",
      "icon": "code-2",
      "desc": "Go 语言相关文章",
      "color": "#3b82f6",
      "count": 24
    }
  ]
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 分类 ID |
| name | string | 分类名称 |
| slug | string | URL 别名 |
| icon | string | lucide 图标名 |
| desc | string | 分类描述（对应表字段 description） |
| color | string | 主题色（十六进制） |
| count | int | 该分类下已发布文章数量 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 500 | 500 | 获取分类列表失败 |

---

#### 10. 前台标签列表

请求方式：`GET /api/v1/tags`

认证方式：无

接口说明：返回所有启用（`status = active`）的标签及使用次数，按 ID 升序。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": [
    { "id": 1, "name": "Go", "count": 24 }
  ]
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 标签 ID |
| name | string | 标签名称 |
| count | int | 该标签下已发布文章数量 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 500 | 500 | 获取标签列表失败 |

---

#### 11. 前台归档

请求方式：`GET /api/v1/archives`

认证方式：无

接口说明：按「年-月」分组的已发布文章归档列表（含各月文章数）。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": [
    {
      "year": "2026",
      "month": "08",
      "count": 5,
      "articles": [
        {
          "id": 1,
          "slug": "go-concurrency",
          "title": "Go 并发模型",
          "date": "2026-08-18",
          "category": "Go",
          "views": 12840
        }
      ]
    }
  ]
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| year | string | 年份（YYYY） |
| month | string | 月份（MM，两位） |
| count | int | 该月份文章数 |
| articles | array | 该月份文章列表 |
| articles[].id | int | 文章 ID |
| articles[].slug | string | URL 别名 |
| articles[].title | string | 标题 |
| articles[].date | string | 发布日期（YYYY-MM-DD，**[无]** 具体来源格式，由表 `published_at` 派生） |
| articles[].category | string | 分类名称 |
| articles[].views | int | 浏览量 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 500 | 500 | 获取归档列表失败 |

---

#### 12. 前台站点统计

请求方式：`GET /api/v1/stats`

认证方式：无

接口说明：返回站点整体统计数据（文章总数、总浏览量、总评论数、动态列表）。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "articles": 128,
    "views": 1286420,
    "comments": 892,
    "dynamics": [
      { "type": "", "text": "", "time": "", "link": "" }
    ]
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| articles | int | 已发布文章总数 |
| views | int | 已发布文章总浏览量 |
| comments | int | 评论总数 |
| dynamics | array | 动态列表（源码当前逻辑固定返回空数组，**无填充数据逻辑**，即恒为 `[]`） |
| dynamics[].type | string | 动态类型（**[无]** 无实际数据来源） |
| dynamics[].text | string | 动态文本（**[无]**） |
| dynamics[].time | string | 动态时间（**[无]**） |
| dynamics[].link | string | 关联链接（可选字段，**[无]**） |

> 说明：DTO 定义了 `dynamics` 结构，但 `StatsRepository.GetStats` 源码中将其初始化为空 `[]dto.Dynamic{}` 且无赋值逻辑，因此该数组目前恒为空；仅「架构兼容前端展示」使用。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 500 | 500 | 获取站点统计失败 |

---

### 评论

#### 13. 获取文章评论列表

请求方式：`GET /api/v1/comments/:slug`

认证方式：无

接口说明：分页获取某篇文章的一级评论（`status = normal` 且 `parent_id IS NULL`），按创建时间倒序。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| slug | string | 是 | path | 文章 URL 别名 |
| page | int | 否 | query | 页码，默认 1 |
| page_size | int | 否 | query | 每页条数，默认 10（上限 50） |

请求示例：

```
GET /api/v1/comments/go-concurrency?page=1&page_size=10
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "parent_id": null,
        "user_name": "Kevin_z",
        "avatar": "https://...",
        "content": "写得很清楚",
        "time": "2026-08-22T14:32:00+08:00",
        "like_count": 12,
        "liked": false,
        "reply_total": 3,
        "has_more_reply": false
      }
    ],
    "total": 56
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| list | array | 评论列表 |
| list[].id | int | 评论 ID |
| list[].parent_id | int/null | 父评论 ID（一级评论为 null） |
| list[].user_name | string | 评论者昵称 |
| list[].avatar | string | 评论者头像 URL |
| list[].content | string | 评论内容 |
| list[].time | string | 评论时间（RFC3339） |
| list[].like_count | int | 点赞数 |
| list[].liked | bool | 当前用户是否已点赞（本接口未登录恒为 false，**[无]** 登录态点赞状态判断逻辑） |
| list[].replies | array | 内嵌回复列表（`omitempty`，本列表接口**不返回该字段**，回复需走「回复列表」接口） |
| list[].reply_total | int | 子评论总数 |
| list[].has_more_reply | bool | 是否还有更多回复 |
| total | int | 评论总数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 参数错误：分页参数绑定失败 |
| 500 | 500 | 获取评论列表失败（含 slug 对应文章不存在的情形，message 统一为该值） |

> 说明：DTO 中定义了 `liked`、`replies`（`omitempty`）字段；本前台评论列表查询未携带登录用户，`liked` 固定为 `false`，且 `replies` 不会被填充（`omitempty` 故不输出）。`liked` 的登录态判断在源码中未明确，标注 **[无]**。

---

#### 14. 获取评论的回复列表

请求方式：`GET /api/v1/comments/replies/:id`

认证方式：无

接口说明：分页获取某条评论（或回复）下的回复列表。按创建时间倒序。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 评论 ID（被回复的父评论 ID） |
| page | int | 否 | query | 页码，默认 1 |
| page_size | int | 否 | query | 每页条数，默认 10（上限 50） |

请求示例：

```
GET /api/v1/comments/replies/1?page=1&page_size=10
```

成功响应：结构与「获取文章评论列表」的 `data` 相同（`list` / `total`），列表元素字段同表 13。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 评论 ID 格式错误（路径非数字）/ 参数错误（分页绑定失败） |
| 500 | 500 | 获取回复列表失败 |

---

#### 15. 发表评论

请求方式：`POST /api/v1/comments/:slug`

认证方式：user

接口说明：登录用户对某篇文章发表评论或回复。评论人信息（昵称、头像、IP）由服务端从令牌与库中获取，前端**无需且不可**传昵称/头像。新评论默认状态 `normal`。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| slug | string | 是 | path | 文章 URL 别名 |
| content | string | 是 | body | 评论内容，1~1000 字 |
| parent_id | int/null | 否 | body | 父评论 ID；为空表示一级评论，非空表示回复 |

请求示例：

```json
{
  "content": "写得很清楚，学习了！",
  "parent_id": null
}
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "id": 9,
    "status": "normal",
    "message": "评论已提交"
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 新建评论 ID |
| status | string | 固定 `normal` |
| message | string | 固定文案「评论已提交」 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 401 | 401 | 请先登录 / 认证失败（未取到上下文用户 ID 或令牌校验失败） |
| 400 | 400 | 参数错误（JSON 绑定失败或 content 长度/必填不满足）/ 获取头像失败 |
| 500 | 500 | 发表评论失败：…（内部错误，如文章不存在、写入失败，message 带前缀「发表评论失败：」） |

---

#### 16. 文章点赞 / 取消点赞

请求方式：`POST /api/v1/like/article/:slug`

认证方式：user

接口说明：切换当前用户对某篇文章的点赞状态——未点赞则点赞，已点赞则取消（幂等切换）。返回最新点赞状态与点赞数。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| slug | string | 是 | path | 文章 URL 别名 |

请求示例：

```
POST /api/v1/like/article/go-concurrency
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "liked": true,
    "likes": 863
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| liked | bool | 本次操作后是否已点赞 |
| likes | int | 文章最新点赞总数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 401 | 401 | 请先登录 / 认证失败 |
| 500 | 500 | 点赞失败 |

---

#### 17. 评论点赞 / 取消点赞

请求方式：`POST /api/v1/like/comment/:id`

认证方式：user

接口说明：切换当前用户对某条评论的点赞状态（幂等切换）。返回最新状态与点赞数。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 评论 ID |

请求示例：

```
POST /api/v1/like/comment/9
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "liked": true,
    "like_count": 13
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| liked | bool | 本次操作后是否已点赞 |
| like_count | int | 评论最新点赞总数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 评论 ID 格式错误（路径非数字） |
| 401 | 401 | 请先登录 / 认证失败 |
| 500 | 500 | 点赞失败 |

---

### 后台认证

#### 18. 获取图形验证码

请求方式：`GET /api/v1/admin/auth/captcha`

认证方式：无

接口说明：获取后台登录所需的图形验证码。返回验证码 ID 与 base64 图片数据。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "captcha_id": "<uuid>",
    "captcha_image": "data:image/png;base64,..."
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| captcha_id | string | 验证码唯一 ID，登录时需原样回传 |
| captcha_image | string | 验证码图片（base64 数据），可直接用作 `<img src>` |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 50000 | 500 | 生成验证码失败 |

---

#### 19. 后台管理员登录

请求方式：`POST /api/v1/admin/auth/login`

认证方式：无

接口说明：管理员账号密码 + 图形验证码登录。成功返回 token 对与管理员信息，并执行单设备登录（旧令牌作废）。

请求参数（Body，JSON）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| username | string | 是 | body | 管理员登录名 |
| password | string | 是 | body | 密码（至少 6 位） |
| captcha_id | string | 是 | body | 验证码 ID（来自验证码接口） |
| captcha | string | 是 | body | 用户输入的验证码 |

请求示例：

```json
{
  "username": "admin",
  "password": "123456",
  "captcha_id": "<uuid>",
  "captcha": "a1b2"
}
```

成功响应：

```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "access_token": "<jwt>",
    "refresh_token": "<jwt>",
    "user": {
      "id": 1,
      "provider": "admin",
      "nickname": "阿轩",
      "avatar": ""
    }
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| access_token | string | 访问令牌 |
| refresh_token | string | 刷新令牌 |
| user | object | 管理员信息 |
| user.id | int | 管理员 ID |
| user.provider | string | 固定 `admin` |
| user.nickname | string | 显示昵称 |
| user.avatar | string | 头像 URL（可为空） |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40000 | 400 | 请求参数错误（JSON 绑定失败） |
| 3004 | 400 | 验证码错误（图形验证码校验不通过） |
| 40100 | 401 | 登录失败：用户名或密码错误（message 前缀「登录失败：」，实际取 1003 默认文案或底层错误） |

---

#### 20. 后台刷新令牌

请求方式：`POST /api/v1/admin/auth/refresh`

认证方式：无

接口说明：使用 refresh_token 换取新的 access_token（管理员侧）。

请求参数（Body，JSON）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| refresh_token | string | 是 | body | 登录时下发的 refresh_token |

请求示例：

```json
{
  "refresh_token": "<jwt>"
}
```

成功响应：

```json
{
  "code": 0,
  "message": "刷新成功",
  "data": {
    "access_token": "<jwt>"
  }
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40000 | 400 | 请求参数错误（JSON 绑定失败） |
| 40100 | 401 | 刷新失败：refresh_token 无效 / 已失效 / 检查黑名单失败 / 生成失败（message 前缀「刷新失败：」） |

---

#### 21. 后台退出登录

请求方式：`POST /api/v1/admin/auth/logout`

认证方式：admin

接口说明：将当前 access_token 加入黑名单并清除活跃令牌。令牌无效/过期仍返回成功（幂等）。

请求参数：无（从 `Authorization` 头解析令牌）

成功响应：

```json
{
  "code": 0,
  "message": "退出成功"
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40000 | 400 | 缺少 access_token（Authorization 头为空） |
| 50000 | 500 | 退出失败（message 前缀「退出失败：」） |

---

#### 22. 获取管理员资料

请求方式：`GET /api/v1/admin/profile`

认证方式：admin

接口说明：返回当前登录管理员的用户名与头像。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "username": "admin",
    "avatar": ""
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| username | string | 管理员用户名 |
| avatar | string | 头像 URL（可为空） |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40100 | 401 | 未登录（未取到上下文 user_id） |
| 50000 | 500 | 获取资料失败（管理员不存在，message 前缀「获取资料失败：」）**[注：源码 GetProfile 失败走 InternalError 分支，HTTP 500]** |

---

#### 23. 更新管理员资料

请求方式：`PUT /api/v1/admin/profile`

认证方式：admin

接口说明：更新当前管理员资料。`username`/`avatar` 非空时更新；`password` 非空时重置密码（bcrypt 加密）。**username 可为空（不修改），无唯一性/长度二次校验（除登录名由表结构唯一索引保证）**。

请求参数（Body，JSON）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| username | string | 否 | body | 新登录名（空则不修改） |
| avatar | string | 否 | body | 新头像 URL（空则不修改） |
| password | string | 否 | body | 新密码（空则不修改密码） |

请求示例：

```json
{
  "username": "admin",
  "avatar": "https://.../avatar.png",
  "password": "newpass123"
}
```

成功响应：

```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "username": "admin",
    "avatar": "https://.../avatar.png"
  }
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40100 | 401 | 未登录 |
| 40000 | 400 | 请求参数错误（JSON 绑定失败） |
| 50000 | 500 | 更新失败：…（管理员不存在 / 密码加密失败 / 保存失败，message 前缀「更新失败：」） |

---

### 后台文章管理

> 说明：文章管理接口均挂载 `OperationLog` 中间件，写操作成功（HTTP 2xx）后自动记录操作日志。列表/详情/详情含全部状态（草稿/已发布/已下架），与前台仅看已发布不同。

#### 24. 后台文章列表

请求方式：`GET /api/v1/admin/articles`

认证方式：admin

接口说明：分页获取全部状态文章，支持关键词、分类名、标签名、状态筛选。

请求参数（Query）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| page | int | 否 | query | 页码，默认 1（min=1） |
| page_size | int | 否 | query | 每页条数，默认 10，上限 100 |
| keyword | string | 否 | query | 标题关键词模糊搜索 |
| category | string | 否 | query | 按**分类名称**筛选（传 `all` 或不传为不过滤） |
| tag | string | 否 | query | 按**标签名称**筛选（传 `all` 或不传为不过滤） |
| status | string | 否 | query | 按状态筛选，取值 `published` / `draft` / `down` / `all` |

请求示例：

```
GET /api/v1/admin/articles?page=1&page_size=10&status=published&category=Go
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "slug": "go-concurrency",
        "title": "Go 并发模型",
        "summary": "摘要…",
        "cover_image": "https://...",
        "category": "Go",
        "tags": ["Go", "并发"],
        "status": "published",
        "views": 12840,
        "date": "2026-08-18",
        "updated": "2026-08-22"
      }
    ],
    "total": 128
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| list | array | 文章列表 |
| list[].id | int | 文章 ID |
| list[].slug | string | URL 别名 |
| list[].title | string | 标题 |
| list[].summary | string | 摘要 |
| list[].cover_image | string | 封面图 URL |
| list[].category | string | 分类名称 |
| list[].tags | array | 标签名称列表 |
| list[].status | string | 状态：`published` / `draft` / `down` |
| list[].views | int | 浏览量 |
| list[].date | string | 发布时间（YYYY-MM-DD） |
| list[].updated | string | 更新时间（YYYY-MM-DD） |
| total | int | 总记录数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 请求参数错误（分页/状态等绑定失败） |
| 500 | 500 | 获取文章列表失败 |

---

#### 25. 后台文章详情

请求方式：`GET /api/v1/admin/articles/:id`

认证方式：admin

接口说明：按 ID 获取文章详情（含全部状态与完整 Markdown 正文）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 文章 ID |

请求示例：

```
GET /api/v1/admin/articles/1
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "id": 1,
    "slug": "go-concurrency",
    "title": "Go 并发模型",
    "summary": "摘要…",
    "cover_image": "https://...",
    "category": "Go",
    "tags": ["Go", "并发"],
    "status": "published",
    "views": 12840,
    "date": "2026-08-18",
    "updated": "2026-08-22",
    "content": "# Markdown 正文…"
  }
}
```

响应字段说明：在「文章列表项」基础上增加：

| 字段 | 类型 | 说明 |
|------|------|------|
| content | string | Markdown 正文 |

（其余字段同接口 24 列表项。）

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 文章 ID 格式错误（路径非数字） |
| 404 | 404 | 文章不存在 |

---

#### 26. 新建文章

请求方式：`POST /api/v1/admin/articles`

认证方式：admin

接口说明：新建文章。分类传**名称**，标签传**名称数组**（后端按名称查找/创建分类标签并建立关联）。作者取当前登录管理员。

请求参数（Body，JSON）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| title | string | 是 | body | 标题（≤200） |
| summary | string | 否 | body | 摘要（≤500） |
| cover_image | string | 否 | body | 封面图 URL |
| category | string | 是 | body | 分类名称（≤50） |
| tags | array | 否 | body | 标签名称列表 |
| content | string | 是 | body | Markdown 正文 |
| status | string | 是 | body | `published` / `draft` / `down` |
| published_at | string | 否 | body | 发布时间，格式 `YYYY-MM-DD` |

请求示例：

```json
{
  "title": "Go 并发模型",
  "summary": "从 goroutine 到 channel",
  "cover_image": "https://.../cover.png",
  "category": "Go",
  "tags": ["Go", "并发"],
  "content": "# 正文",
  "status": "published",
  "published_at": "2026-08-18"
}
```

成功响应：`data` 为该新文章详情（结构同接口 25）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 请求参数错误（JSON 绑定失败，title/category/content/status 缺失或非法） |
| 401 | 401 | 未登录（上下文无作者 ID） |
| 500 | 500 | 创建文章失败：…（内部错误，message 带底层错误） |

---

#### 27. 更新文章

请求方式：`PUT /api/v1/admin/articles/:id`

认证方式：admin

接口说明：按 ID 全量更新文章（title/category/content/status 等为必填，接口以整篇覆盖语义提交）。分类/标签同上按名称处理。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 文章 ID |
| body | object | 是 | body | 字段同「新建文章」请求体（title/summary/cover_image/category/tags/content/status/published_at） |

请求示例：

```json
{
  "title": "Go 并发模型（修订）",
  "summary": "更新后的摘要",
  "cover_image": "https://.../cover.png",
  "category": "Go",
  "tags": ["Go", "并发", "面试"],
  "content": "# 更新后的正文",
  "status": "published",
  "published_at": "2026-08-22"
}
```

成功响应：`data` 为更新后的文章详情（结构同接口 25）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 文章 ID 格式错误 / 请求参数错误（JSON 绑定失败） |
| 500 | 500 | 更新文章失败：…（内部错误，message 带底层错误） |

---

#### 28. 删除文章

请求方式：`DELETE /api/v1/admin/articles/:id`

认证方式：admin

接口说明：软删除指定文章（写 `deleted_at`）。源码仓储层对不存在记录的删除结果处理未明确，**[无]** 明确「找不到」分支。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 文章 ID |

成功响应：

```json
{
  "code": 0,
  "message": "成功"
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 文章 ID 格式错误 |
| 500 | 500 | 删除文章失败 |

---

#### 29. 修改文章状态

请求方式：`PUT /api/v1/admin/articles/:id/status`

认证方式：admin

接口说明：发布 / 存草稿 / 下架文章。操作日志中间件将据 body 的 `status` 推断动作（`published`→发布、`draft`→存为草稿、`down`→下架）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 文章 ID |
| status | string | 是 | body | `published` / `draft` / `down` |

请求示例：

```json
{
  "status": "published"
}
```

成功响应：`data` 为更新后的文章详情（结构同接口 25）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 文章 ID 格式错误 / 请求参数错误（JSON 绑定失败或 status 非法） |
| 500 | 500 | 内部错误（更新失败时返回，message 带底层错误） |

---

### 后台分类管理

> 说明：分类管理接口均挂载 `OperationLog` 中间件。

#### 30. 后台分类列表

请求方式：`GET /api/v1/admin/categories`

认证方式：admin

接口说明：分页获取分类列表，含文章数统计，支持名称、状态、文章数区间筛选。

请求参数（Query）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| page | int | 否 | query | 页码，默认 1 |
| page_size | int | 否 | query | 每页条数，默认 10，上限 100 |
| keyword | string | 否 | query | 名称模糊搜索 |
| status | string | 否 | query | `active` / `inactive` / `all`（默认不过滤） |
| min_count | int | 否 | query | 最小文章数 |
| max_count | int | 否 | query | 最大文章数 |

请求示例：

```
GET /api/v1/admin/categories?page=1&page_size=10&status=active
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "Go",
        "slug": "go",
        "icon": "code-2",
        "desc": "Go 语言相关",
        "color": "#3b82f6",
        "status": "active",
        "count": 24,
        "created_at": "2024-01-05",
        "updated_at": "2026-08-20"
      }
    ],
    "total": 11
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| list[].id | int | 分类 ID |
| list[].name | string | 名称 |
| list[].slug | string | URL 别名 |
| list[].icon | string | 图标名 |
| list[].desc | string | 描述（表字段 description） |
| list[].color | string | 主题色 |
| list[].status | string | `active` / `inactive` |
| list[].count | int | 已发布文章数 |
| list[].created_at | string | 创建时间（YYYY-MM-DD） |
| list[].updated_at | string | 更新时间（YYYY-MM-DD） |
| total | int | 总记录数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 请求参数错误 |
| 500 | 500 | 获取分类列表失败 |

---

#### 31. 后台分类详情

请求方式：`GET /api/v1/admin/categories/:id`

认证方式：admin

接口说明：按 ID 获取分类详情（单条结构同接口 30 列表项）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 分类 ID |

成功响应：`data` 为单个分类对象（结构同接口 30 列表项字段）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 分类 ID 格式错误 |
| 404 | 404 | 分类不存在 |

---

#### 32. 新建分类

请求方式：`POST /api/v1/admin/categories`

认证方式：admin

接口说明：新建分类。`slug` 留空则自动生成；`status` 缺省 `active`。名称唯一冲突会返回错误。

请求参数（Body，JSON）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| name | string | 是 | body | 名称（≤50） |
| slug | string | 否 | body | URL 别名（≤60，留空自动生成） |
| desc | string | 否 | body | 描述（≤255） |
| color | string | 否 | body | 主题色（≤10） |
| status | string | 否 | body | `active` / `inactive`，缺省 `active` |

请求示例：

```json
{
  "name": "MySQL",
  "slug": "mysql",
  "desc": "数据库相关",
  "color": "#22c55e",
  "status": "active"
}
```

成功响应：`data` 为新建分类对象（结构同接口 30 列表项字段）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 请求参数错误（JSON 绑定失败）/ 分类名称已存在（唯一约束冲突） |
| 500 | 500 | 创建分类失败 |

---

#### 33. 更新分类

请求方式：`PUT /api/v1/admin/categories/:id`

认证方式：admin

接口说明：部分更新分类（传哪些字段更新哪些，字段均为可选）。名称冲突返回错误。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 分类 ID |
| name | string | 否 | body | 名称（≤50） |
| slug | string | 否 | body | URL 别名（≤60） |
| desc | string | 否 | body | 描述（≤255） |
| color | string | 否 | body | 主题色（≤10） |
| status | string | 否 | body | `active` / `inactive` |

请求示例：

```json
{
  "name": "MySQL 数据库",
  "color": "#16a34a"
}
```

成功响应：`data` 为更新后的分类对象。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 分类 ID 格式错误 / 请求参数错误 / 分类名称已存在 |
| 500 | 500 | 更新分类失败 |

---

#### 34. 修改分类状态

请求方式：`PUT /api/v1/admin/categories/:id/status`

认证方式：admin

接口说明：启用 / 禁用分类（操作日志推断：`active`→启用、`inactive`→禁用）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 分类 ID |
| status | string | 是 | body | `active` / `inactive` |

请求示例：

```json
{
  "status": "inactive"
}
```

成功响应：

```json
{
  "code": 0,
  "message": "成功"
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 分类 ID 格式错误 / 请求参数错误（status 非法） |
| 500 | 500 | 修改状态失败 |

---

#### 35. 删除分类

请求方式：`DELETE /api/v1/admin/categories/:id`

认证方式：admin

接口说明：删除指定分类。分类下存在关联文章时的行为源码未明确（**[无]** 级联或限制逻辑）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 分类 ID |

成功响应：

```json
{
  "code": 0,
  "message": "成功"
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 分类 ID 格式错误 |
| 500 | 500 | 删除分类失败 |

---

### 后台标签管理

> 说明：标签管理接口均挂载 `OperationLog` 中间件。

#### 36. 后台标签列表

请求方式：`GET /api/v1/admin/tags`

认证方式：admin

接口说明：分页获取标签列表，含使用次数统计，支持名称、状态、使用次数区间筛选。

请求参数（Query）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| page | int | 否 | query | 页码，默认 1 |
| page_size | int | 否 | query | 每页条数，默认 10，上限 100 |
| keyword | string | 否 | query | 名称模糊搜索 |
| status | string | 否 | query | `active` / `inactive` / `all`（默认不过滤） |
| min_count | int | 否 | query | 最小使用次数 |
| max_count | int | 否 | query | 最大使用次数 |

请求示例：

```
GET /api/v1/admin/tags?page=1&page_size=10
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "Go",
        "status": "active",
        "count": 24,
        "created_at": "2024-01-05",
        "updated_at": "2026-08-20"
      }
    ],
    "total": 46
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| list[].id | int | 标签 ID |
| list[].name | string | 标签名称 |
| list[].status | string | `active` / `inactive` |
| list[].count | int | 关联已发布文章数 |
| list[].created_at | string | 创建时间（YYYY-MM-DD） |
| list[].updated_at | string | 更新时间（YYYY-MM-DD） |
| total | int | 总记录数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 请求参数错误 |
| 500 | 500 | 获取标签列表失败 |

---

#### 37. 后台标签详情

请求方式：`GET /api/v1/admin/tags/:id`

认证方式：admin

接口说明：按 ID 获取标签详情（单条结构同接口 36 列表项）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 标签 ID |

成功响应：`data` 为单个标签对象（结构同接口 36 列表项字段）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 标签 ID 格式错误 |
| 404 | 404 | 标签不存在 |

---

#### 38. 新建标签

请求方式：`POST /api/v1/admin/tags`

认证方式：admin

接口说明：新建标签。`status` 缺省 `active`；名称唯一冲突返回错误。

请求参数（Body，JSON）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| name | string | 是 | body | 标签名称（≤50） |
| status | string | 否 | body | `active` / `inactive`，缺省 `active` |

请求示例：

```json
{
  "name": "分布式",
  "status": "active"
}
```

成功响应：`data` 为新建标签对象（结构同接口 36 列表项字段）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 请求参数错误（JSON 绑定失败）/ 标签名称已存在 |
| 500 | 500 | 创建标签失败 |

---

#### 39. 更新标签

请求方式：`PUT /api/v1/admin/tags/:id`

认证方式：admin

接口说明：部分更新标签（仅 `name` 生效；`status` 字段虽可传但 service 未处理——源码中 `AdminTagUpdateRequest` 含 status，但 `TagsAdminService.AdminUpdate` 只更新 name，**[源码已核实]**）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 标签 ID |
| name | string | 否 | body | 新名称（≤50） |

请求示例：

```json
{
  "name": "分布式系统"
}
```

成功响应：`data` 为更新后的标签对象。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 标签 ID 格式错误 / 请求参数错误 / 标签名称已存在 |
| 500 | 500 | 更新标签失败 |

---

#### 40. 修改标签状态

请求方式：`PUT /api/v1/admin/tags/:id/status`

认证方式：admin

接口说明：启用 / 禁用标签（操作日志推断：`active`→启用、`inactive`→禁用）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 标签 ID |
| status | string | 是 | body | `active` / `inactive` |

请求示例：

```json
{
  "status": "inactive"
}
```

成功响应：

```json
{
  "code": 0,
  "message": "成功"
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 标签 ID 格式错误 / 请求参数错误（status 非法） |
| 500 | 500 | 修改状态失败 |

---

#### 41. 删除标签

请求方式：`DELETE /api/v1/admin/tags/:id`

认证方式：admin

接口说明：删除指定标签。存在关联文章时的行为源码未明确（**[无]** 级联或限制逻辑）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 标签 ID |

成功响应：

```json
{
  "code": 0,
  "message": "成功"
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 标签 ID 格式错误 |
| 500 | 500 | 删除标签失败 |

---

### 后台评论管理

> 说明：评论管理接口均挂载 `OperationLog` 中间件。

#### 42. 后台评论列表

请求方式：`GET /api/v1/admin/comments`

认证方式：admin

接口说明：分页获取评论（含全部状态），支持内容/用户名关键词、状态、文章、时间段筛选。

请求参数（Query）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| page | int | 否 | query | 页码，默认 1 |
| page_size | int | 否 | query | 每页条数，默认 10，上限 100 |
| keyword | string | 否 | query | 内容或用户名模糊搜索 |
| status | string | 否 | query | `normal` / `banned` / `all`（默认不过滤） |
| article_id | int | 否 | query | 按文章 ID 筛选 |
| start_date | string | 否 | query | 开始日期（`YYYY-MM-DD`），按创建时间 |
| end_date | string | 否 | query | 结束日期（`YYYY-MM-DD`），按创建时间 |

请求示例：

```
GET /api/v1/admin/comments?page=1&page_size=10&status=banned
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 4,
        "article_id": 1,
        "article_title": "Go 并发模型",
        "parent_id": null,
        "user_id": 1,
        "user_name": "Kevin_z",
        "user_avatar": "https://...",
        "content": "内容…",
        "status": "banned",
        "like_count": 3,
        "ip": "39.156.***.***",
        "created_at": "2026-07-29 10:02"
      }
    ],
    "total": 2
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| list[].id | int | 评论 ID |
| list[].article_id | int | 所属文章 ID |
| list[].article_title | string | 所属文章标题 |
| list[].parent_id | int/null | 父评论 ID |
| list[].user_id | int | 用户 ID |
| list[].user_name | string | 用户昵称 |
| list[].user_avatar | string | 用户头像 |
| list[].content | string | 评论内容 |
| list[].status | string | `normal` / `banned` |
| list[].like_count | int | 点赞数 |
| list[].ip | string | 评论者 IP |
| list[].created_at | string | 创建时间（格式见上示例，**[无]** 精确布局由格式化确定） |
| total | int | 总记录数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 请求参数错误 |
| 500 | 500 | 获取评论列表失败 |

---

#### 43. 后台评论详情

请求方式：`GET /api/v1/admin/comments/:id`

认证方式：admin

接口说明：按 ID 获取评论详情（单条结构同接口 42 列表项）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 评论 ID |

成功响应：`data` 为单个评论对象（结构同接口 42 列表项字段）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 评论 ID 格式错误 |
| 404 | 404 | 评论不存在 |

---

#### 44. 修改评论状态

请求方式：`PUT /api/v1/admin/comments/:id/status`

认证方式：admin

接口说明：解封 / 封禁评论（操作日志推断：`normal`→解封、`banned`→封禁）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 评论 ID |
| status | string | 是 | body | `normal` / `banned` |

请求示例：

```json
{
  "status": "banned"
}
```

成功响应：`data` 为更新后的评论对象（结构同接口 42 列表项字段）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 评论 ID 格式错误 / 请求参数错误（status 非法） |
| 500 | 500 | 修改评论状态失败 |

---

#### 45. 删除评论

请求方式：`DELETE /api/v1/admin/comments/:id`

认证方式：admin

接口说明：软删除指定评论。存在回复时的处理逻辑源码未明确（**[无]**）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 评论 ID |

成功响应：

```json
{
  "code": 0,
  "message": "成功"
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 评论 ID 格式错误 |
| 500 | 500 | 删除评论失败 |

---

### 后台用户管理

> 说明：用户管理接口均挂载 `OperationLog` 中间件。管理对象为前台注册用户（GitHub OAuth 用户），非管理员。

#### 46. 后台用户列表

请求方式：`GET /api/v1/admin/users`

认证方式：admin

接口说明：分页获取前台用户列表，支持昵称、状态、注册时间段筛选。

请求参数（Query）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| page | int | 否 | query | 页码，默认 1 |
| page_size | int | 否 | query | 每页条数，默认 10，上限 100 |
| keyword | string | 否 | query | 昵称模糊搜索 |
| status | string | 否 | query | `active`（status=1）/ `inactive`（status=0）/ `all`（默认不过滤） |
| start_date | string | 否 | query | 开始日期（`YYYY-MM-DD`），按创建时间 |
| end_date | string | 否 | query | 结束日期（`YYYY-MM-DD`），按创建时间 |

请求示例：

```
GET /api/v1/admin/users?page=1&page_size=10&status=active
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "provider": "github",
        "openid": "12345678",
        "nickname": "Kevin_z",
        "avatar": "https://...",
        "status": 1,
        "last_login_at": "2026-08-22 10:30",
        "created_at": "2026-07-01 12:00"
      }
    ],
    "total": 128
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| list[].id | int | 用户 ID |
| list[].provider | string | 登录方式（github） |
| list[].openid | string | 平台用户标识 |
| list[].nickname | string | 昵称 |
| list[].avatar | string | 头像 URL |
| list[].status | int | 状态：`1` 正常 / `0` 禁用 |
| list[].last_login_at | string | 最后登录时间 |
| list[].created_at | string | 注册时间 |
| total | int | 总记录数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 请求参数错误 |
| 500 | 500 | 获取用户列表失败 |

---

#### 47. 后台用户详情

请求方式：`GET /api/v1/admin/users/:id`

认证方式：admin

接口说明：按 ID 获取前台用户详情（单条结构同接口 46 列表项）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 用户 ID |

成功响应：`data` 为单个用户对象（结构同接口 46 列表项字段）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 用户 ID 格式错误 |
| 404 | 404 | 用户不存在 |

---

#### 48. 修改用户状态

请求方式：`PUT /api/v1/admin/users/:id/status`

认证方式：admin

接口说明：启用 / 禁用前台用户。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 用户 ID |
| status | int | 是 | body | `1` 正常 / `0` 禁用 |

请求示例：

```json
{
  "status": 0
}
```

成功响应：`data` 为更新后的用户对象（结构同接口 46 列表项字段）。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 用户 ID 格式错误 / 请求参数错误（status 非 0/1） |
| 500 | 500 | 修改用户状态失败 |

---

#### 49. 删除用户

请求方式：`DELETE /api/v1/admin/users/:id`

认证方式：admin

接口说明：删除指定前台用户。用户存在评论/点赞等关联数据时的处理源码未明确（**[无]**）。

请求参数：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| id | int | 是 | path | 用户 ID |

成功响应：

```json
{
  "code": 0,
  "message": "成功"
}
```

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 用户 ID 格式错误 |
| 500 | 500 | 删除用户失败 |

---

### 系统设置

> 说明：系统设置接口挂载 `OperationLog` 中间件。`site_settings` 为键值对表，源码未内置完整的键清单，可用的键由「前台站点信息」各字段对应的配置键推导（如 `blog_name`、`tagline`、`blog_desc`、`author_intro`、`author_name`、`author_role`、`github`、`email`、`motto`、`location`、`since`、`avatar`、`socials`、`hero_terminal` 等）；其余键是否生效取决于前端使用方，后端仅作 KV 原样存取。

#### 50. 读取系统设置

请求方式：`GET /api/v1/admin/settings`

认证方式：admin

接口说明：返回全部站点设置项（key → value）。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "blog_name": "小猫的个人博客",
    "tagline": "记录代码，分享技术",
    "author_name": "阿轩",
    "since": "2019"
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| data | object | 全部设置项，键值均为字符串（`data` 为任意 key 的 map，key/value 均为 string） |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 50000 | 500 | 获取设置失败 |

---

#### 51. 保存系统设置

请求方式：`PUT /api/v1/admin/settings`

认证方式：admin

接口说明：批量保存设置项（逐条 upsert：键存在则更新，不存在则插入）。成功后返回更新后的完整设置。**注意：本接口请求体为任意字符串键值 map，不对具体键做白名单校验（源码未明确键约束）。**

请求参数（Body，JSON）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| (任意键) | string | 是 | body | 形如 `{ "blog_name": "...", ... }` 的设置键值，至少需一个键 |

请求示例：

```json
{
  "blog_name": "小猫的个人博客",
  "tagline": "记录代码，分享技术，持续成长",
  "since": "2019"
}
```

成功响应：`data` 为保存后重新读取的完整设置对象（同接口 50 的 `data` 结构）。若保存成功但随后重读失败，则返回的 `data` 为请求提交的 `body` 原样。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 40001 | 400 | 请求参数错误（JSON 绑定失败，如非对象或值非字符串） |
| 50000 | 500 | 保存设置失败 |

---

### 文件上传

#### 52. 上传图片

请求方式：`POST /api/v1/upload/image`

认证方式：admin

接口说明：上传单张图片，返回可访问 URL。后端会做大小与格式校验，并对 jpg/png 去除 EXIF 与等比压缩（最长边超 1920px 压缩），gif/webp 原样上传。存储驱动（OSS/本地）由配置决定。

请求参数（`multipart/form-data`）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| file | file | 是 | form | 图片文件；仅支持 `jpg/jpeg/png/gif/webp`，大小 ≤ 5MB |

请求示例：

```bash
curl -X POST /api/v1/upload/image \
  -H "Authorization: Bearer <access_token>" \
  -F "file=@cover.png"
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "url": "https://bucket.oss-cn-beijing.aliyuncs.com/images/xxxxx.jpg"
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| url | string | 上传后图片可访问 URL（可直接作为封面/正文图片地址） |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 401 | 401 | 未登录（认证失败或上下文无用户 ID） |
| 400 | 400 | 参数错误：请选择要上传的图片 / 图片不能超过 5MB / 仅支持 jpg/jpeg/png/gif/webp 格式 / 图片文件无效 |
| 500 | 500 | 读取图片失败 / 图片上传失败 |

> 说明：文件类型按**文件头**（mimetype）检测而非扩展名；若校验失败以 `pkg/response` 输出，HTTP 400、code 400。

---

### 仪表盘

> 说明：仪表盘接口仅 admin 认证，不记录操作日志。

#### 53. 仪表盘-统计数据

请求方式：`GET /api/v1/admin/dashboard/stats`

认证方式：admin

接口说明：返回文章总数、启用分类数、启用标签数、评论总数的统计（不含浏览量/访问量）。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "articles": { "value": 128 },
    "categories": { "value": 11 },
    "tags": { "value": 46 },
    "comments": { "value": 892 }
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| articles | object | 文章总数 |
| articles.value | int | 数量 |
| categories | object | 启用（active）分类数 |
| tags | object | 启用（active）标签数 |
| comments | object | 评论总数 |
| (各项).value | int | 对应数量 |

> 说明：统计维度源码固定为 `articles/categories/tags/comments` 四个 key（`map[string]StatItem`），仅含 `value`。

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 500 | 500 | 获取统计数据失败 |

---

#### 54. 仪表盘-最新文章

请求方式：`GET /api/v1/admin/dashboard/articles`

认证方式：admin

接口说明：返回最近发布的若干篇文章（源码取最新发布的文章，数量逻辑见仓储层，**[无]** 明确条数上限说明；结构按列表返回）。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "recent_posts": [
      {
        "id": 1,
        "title": "Go 并发模型",
        "category": "Go",
        "date": "2026-08-18",
        "views": 12840
      }
    ]
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| recent_posts | array | 最近发布文章列表（注意外层包在 `data` 下） |
| recent_posts[].id | int | 文章 ID |
| recent_posts[].title | string | 标题 |
| recent_posts[].category | string | 分类名称 |
| recent_posts[].date | string | 发布日期 |
| recent_posts[].views | int | 浏览量 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 500 | 500 | 获取文章数据失败 |

---

#### 55. 仪表盘-最新评论

请求方式：`GET /api/v1/admin/dashboard/comments`

认证方式：admin

接口说明：返回最近若干条评论（源码取最新 5 条）。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": [
    {
      "id": 1,
      "user_name": "Kevin_z",
      "avatar": "https://...",
      "content": "写得好",
      "time": "2026-08-22 14:32"
    }
  ]
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int | 评论 ID |
| user_name | string | 用户名 |
| avatar | string | 头像 URL |
| content | string | 评论内容 |
| time | string | 评论时间 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 500 | 500 | 获取评论数据失败 |

---

#### 56. 仪表盘-最近操作记录

请求方式：`GET /api/v1/admin/dashboard/operations`

认证方式：admin

接口说明：返回最近若干条操作记录（源码取最新 8 条）。

请求参数：无

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": [
    {
      "time": "2026-08-22 10:30",
      "user": "阿轩",
      "action": "发布",
      "target": "文章 #1"
    }
  ]
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| time | string | 操作时间 |
| user | string | 操作者昵称 |
| action | string | 动作（新建/更新/删除/发布/存为草稿/下架/启用/禁用/解封/封禁 等） |
| target | string | 操作对象（如「文章 #1」） |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 500 | 500 | 获取操作记录失败 |

---

### 操作日志

#### 57. 操作日志列表

请求方式：`GET /api/v1/admin/operation-logs`

认证方式：admin

接口说明：分页查询后台操作日志，支持动作、对象、时间段筛选。日志由 `OperationLog` 中间件在后台写操作成功时自动写入。

请求参数（Query）：

| 参数 | 类型 | 必填 | 位置 | 说明 |
|------|------|------|------|------|
| page | int | 否 | query | 页码，默认 1 |
| page_size | int | 否 | query | 每页条数，默认 10，上限 100 |
| action | string | 否 | query | 按动作筛选（新建/更新/删除/发布/存为草稿/下架/启用/禁用/解封/封禁） |
| target | string | 否 | query | 按操作对象模糊搜索 |
| start_date | string | 否 | query | 开始日期（`YYYY-MM-DD`） |
| end_date | string | 否 | query | 结束日期（`YYYY-MM-DD`，含当天） |

请求示例：

```
GET /api/v1/admin/operation-logs?page=1&page_size=10&action=发布
```

成功响应：

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "admin_id": 1,
        "admin_name": "阿轩",
        "action": "发布",
        "target": "文章 #1",
        "created_at": "2026-08-22 10:30"
      }
    ],
    "total": 45
  }
}
```

响应字段说明：

| 字段 | 类型 | 说明 |
|------|------|------|
| list[].id | int | 日志 ID |
| list[].admin_id | int/null | 操作者管理员 ID |
| list[].admin_name | string | 操作者昵称（关联管理员） |
| list[].action | string | 动作 |
| list[].target | string | 操作对象 |
| list[].created_at | string | 操作时间 |
| total | int | 总记录数 |

错误情况：

| code | HTTP | message 说明 |
|------|------|--------------|
| 400 | 400 | 请求参数错误 |
| 500 | 500 | 获取操作日志列表失败 |

---

## 七、分页说明

系统中分页参数、响应结构不统一，前端需按接口区分：

| 接口 | 分页参数 | 响应外层 | 说明 |
|------|----------|----------|------|
| 前台文章列表（#7） | `page` + `size` | `data.list / data.total / data.page / data.size / data.total_page` | 参数必填 |
| 前台评论列表、回复列表（#13/#14） | `page` + `page_size` | `data.list / data.total` | page_size 上限 50，缺省 1/10 |
| 后台各列表（#24/30/36/42/46/57） | `page` + `page_size` | `data.list / data.total` | 默认 1/10，上限 100 |
| 非分页列表（分类/标签/归档/站点统计/仪表盘等） | 无 | `data` 为数组/对象 | — |

通用分页响应（以 `page/size` 类为典型，如前台文章列表）字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| list | array | 当前页数据 |
| total | int | 符合条件的总记录数 |
| page | int | 当前页码 |
| size | int | 每页条数 |
| total_page | int | 总页数（`ceil(total/size)`） |

> 请求分页参数建议范围：`page >= 1`，`page_size ∈ [1, 100]`（评论类为 `[1,50]`）；超范围会触发参数绑定校验失败，返回 400。

---

## 附：源码信息不足项（集中声明）

以下信息在源码中无法确定，文档相应位置已标 **[无]**，前端如需依赖请与后端确认：

1. `GET /comments` 类接口的 `liked` 字段无登录态判断逻辑（固定 false）。
2. `stats` 接口 `dynamics` 恒为空数组，无填充来源。
3. token 过期时长、GitHub OAuth 的 client 配置值属部署配置，运行时由配置文件提供。
4. 各删除类接口（文章/分类/标签/评论/用户）在存在关联数据时的处理（级联/限制/报错文案）未在源码体现。
5. 后台评论 `created_at`、仪表盘时间的精确输出格式，源码仅通过时间格式化生成，具体布局未在本文档锁死（按示例为准）。
6. 归档项 `date` 字段的具体生成格式未在源码中直接标注，按示例 `YYYY-MM-DD` 呈现。

---

*文档结束 · 本文件由源码逐接口核对生成，谨以源码为准。*






