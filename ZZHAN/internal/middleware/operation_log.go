package middleware

import (
	"ZZHAN/internal/model/entity"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// context key：控制器可通过 c.Set(OpActionKey, "禁用") 自定义动作
const OpActionKey = "op_action"

// 资源名映射
var resourceNames = map[string]string{
	"articles":   "文章",
	"categories": "分类",
	"tags":       "标签",
	"comments":   "评论",
	"settings":   "系统设置",
}

// 状态端点动作映射：resource → status → action
var statusActions = map[string]map[string]string{
	"articles": {
		"published": "发布",
		"draft":     "存为草稿",
		"down":      "下架",
	},
	"categories": {
		"active":   "启用",
		"inactive": "禁用",
	},
	"tags": {
		"active":   "启用",
		"inactive": "禁用",
	},
	"comments": {
		"normal": "解封",
		"banned": "封禁",
	},
}

// OperationLog 操作日志中间件，记录后台增删改操作
func OperationLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// 只记录写操作
		if method != http.MethodPost && method != http.MethodPut &&
			method != http.MethodPatch && method != http.MethodDelete {
			c.Next()
			return
		}

		path := c.Request.URL.Path

		// 跳过非业务路径
		if strings.Contains(path, "/auth/") ||
			strings.Contains(path, "/dashboard") ||
			strings.Contains(path, "/profile") ||
			strings.Contains(path, "/upload/") {
			c.Next()
			return
		}

		// 判断是否为状态端点，提前读取 body
		isStatus := strings.HasSuffix(path, "/status")
		var statusVal string
		if isStatus {
			body, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			var payload struct {
				Status string `json:"status"`
			}
			json.Unmarshal(body, &payload)
			statusVal = payload.Status
		}

		// 执行后续 handler
		c.Next()

		// 只在成功时记录（2xx 状态码）
		if c.Writer.Status() >= 400 {
			return
		}

		// 获取管理员 ID
		adminID := int64(GetUserID(c))
		if adminID == 0 {
			return
		}

		// 解析资源类型和 ID
		parts := strings.Split(strings.TrimPrefix(path, "/api/v1/admin/"), "/")
		resourceType := parts[0]
		resourceName := resourceNames[resourceType]
		if resourceName == "" {
			resourceName = resourceType
		}

		// 构造 target
		target := resourceName
		if len(parts) > 1 {
			if id, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
				target = resourceName + " #" + strconv.FormatInt(id, 10)
			}
		}

		// 确定动作（优先级：控制器设置 > 状态端点推断 > HTTP 方法映射）
		var action string

		// 1. 控制器通过 c.Set(OpActionKey, "xxx") 自定义
		if v, ok := c.Get(OpActionKey); ok {
			action, _ = v.(string)
		}

		// 2. 状态端点：根据 resource + status 推断
		if action == "" && isStatus && statusVal != "" {
			if rm, ok := statusActions[resourceType]; ok {
				action = rm[statusVal]
			}
		}

		// 3. HTTP 方法兜底
		if action == "" {
			switch method {
			case http.MethodPost:
				action = "新建"
			case http.MethodPut, http.MethodPatch:
				action = "更新"
			case http.MethodDelete:
				action = "删除"
			}
		}

		// 写入操作日志
		log := entity.OperationLog{
			AdminID: &adminID,
			Action:  action,
			Target:  target,
		}
		db.Create(&log)
	}
}
