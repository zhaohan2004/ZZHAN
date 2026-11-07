# 博客部署

推送或合并到 `main` 后，`Publish Docker images` 工作流构建前后端镜像、同步 MySQL/Redis 镜像，然后通过 SSH 部署。也可以在 Actions 中选择 `main` 手动运行。镜像包需要公开。

## 服务器首次配置

服务器需要 Docker、Compose、curl、flock，且登录账号能够运行 Docker。

`/home/ubuntu/zzhan-blog` 保存：

- `compose.yaml`：本项目运行配置。
- `.env`：按 `.env.example` 填写，密码随机生成，权限 600。实际文件不提交。
- `config.yaml`：从 `config.production.example.yaml` 创建，填写 OSS 等配置，实际文件不提交。

项目目录权限为 700。`config.yaml` 为 644，供容器内非 root 应用用户读取；敏感配置不应放到其他可公开访问的位置。

GitHub Actions Secrets：`SERVER_HOST`、`SERVER_USER`、`SERVER_PORT`、`SERVER_SSH_KEY`、`SERVER_KNOWN_HOSTS`。

MySQL、Redis 不发布宿主机端口，前端使用 3000/4000。数据卷均由 `zzhan-blog` 项目管理，不共用其他项目的数据卷。

## 每次发布

1. 构建镜像，标签为 `sha-完整提交编号`。
2. 上传 Compose 和部署脚本到服务器 `~/zzhan-blog/releases/标签/`。
3. 部署脚本保留现有 `.env` 和 `config.yaml`，只拉取新应用镜像。
4. 更新本项目，等待健康检查，并检查两个端口的页面与 API。
5. 成功后更新 `.env` 中 `APP_TAG`、根目录 Compose 和 `.last-successful-release`。

`DEPENDENCY_TAG` 保持固定，日常发布不会自动升级 MySQL/Redis。依赖镜像仍会同步保存，供以后有计划地升级。

部署会自动执行应用现有的数据库迁移，因此失败后不自动回滚数据库。更新期间可能短暂中断；发布前应备份重要数据库数据。脚本没有删除容器卷、清理全机镜像或停止其他项目的操作。

## 将历史 Base64 文章封面迁移到 OSS

发布新版后，在部署目录使用当前后端配置运行迁移：

```bash
docker compose exec backend /app/server -c /app/config.yaml -migrate-covers
```

迁移逐张上传并验证图片链接可读，再更新对应文章。失败的封面保留原数据，修复原因后可重跑。成功后打开前台首页，在开发者工具的「网络」中确认 `/api/v1/articles` 的封面字段为链接，且响应体积下降。

## 排查

在服务器执行：

```bash
cd ~/zzhan-blog
docker compose ps --all
docker compose logs --tail=100 backend
```

查看日志时请避免把包含敏感信息的输出公开。健康检查只证明应用启动与 HTTP 可访问，OSS 上传、GitHub OAuth 和管理账号登录需要另行验证。

首次上线后仍需完成：更换已暴露的 OSS 密钥、验证上传、配置 GitHub OAuth、建立管理账号及数据库备份。当前 IP + HTTP 访问不提供 HTTPS 加密。
