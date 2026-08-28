package app

import (
	"ZZHAN/internal/api/admin"
	"ZZHAN/internal/api/web"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/storage"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ZZHAN/internal/api"
	"ZZHAN/internal/model/entity"
	"ZZHAN/pkg/config"
	"ZZHAN/pkg/database"
	"ZZHAN/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用结构体
type App struct {
	cfg     *config.Config
	mysqlDB *gorm.DB
	redis   *redis.Client
	router  *api.Router
	server  *http.Server
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{}
}

// Initialize 初始化应用
func (a *App) Initialize(configPath string) error {
	// 1. 加载配置
	if err := a.initConfig(configPath); err != nil {
		return err
	}

	// 2. 初始化日志
	if err := a.initLogger(); err != nil {
		return err
	}

	// 3. 初始化数据库
	if err := a.initDatabase(); err != nil {
		return err
	}

	// 4. 初始化依赖
	a.initDependencies()

	// 5. 初始化路由
	a.initRouter()

	// 6. 初始化服务器
	a.initServer()

	return nil
}

// initConfig 加载配置
func (a *App) initConfig(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	a.cfg = cfg
	return nil
}

// initLogger 初始化日志
func (a *App) initLogger() error {
	if err := logger.Init(&a.cfg.Log); err != nil {
		return fmt.Errorf("日志初始化失败: %w", err)
	}

	// 打印启动横幅
	logger.Info("=========================================")
	logger.Info(fmt.Sprintf("欢迎使用 %s", a.cfg.App.Name))
	logger.Info(fmt.Sprintf("版本: %s", a.cfg.App.Version))
	logger.Info(fmt.Sprintf("模式: %s", a.cfg.App.Mode))
	logger.Info("配置加载成功")
	logger.Info("=========================================")

	return nil
}

// initDatabase 初始化数据库
func (a *App) initDatabase() error {
	// 初始化 MySQL
	mysqlDB, err := database.InitMySQL(&a.cfg.Database.MySQL)
	if err != nil {
		return fmt.Errorf("MySQL 初始化失败: %w", err)
	}
	a.mysqlDB = mysqlDB

	// 自动迁移数据库表
	logger.Info("开始数据库迁移...")
	if err := a.mysqlDB.AutoMigrate(
		// 账号模块（2张）
		&entity.Admin{},
		&entity.User{},
		// 内容模块（4张）
		&entity.Category{},
		&entity.Tag{},
		&entity.Article{},
		&entity.ArticleTag{},
		// 互动模块（3张）
		&entity.Comment{},
		&entity.Like{},
		&entity.CommentLike{},
		// 关于与站点模块（2张）
		&entity.SiteSetting{},
		&entity.AboutItem{},
		// 统计与审计模块（2张）
		&entity.VisitLog{},
		&entity.OperationLog{},
	); err != nil {
		logger.Warn("数据库迁移警告", zap.Error(err))
	} else {
		logger.Info("数据库迁移完成")
	}

	// 初始化 Redis（可选，失败不影响核心功能）
	rs, err := database.InitRedis(&a.cfg.Database.Redis)
	if err != nil {
		logger.Warn("Redis 初始化失败，将不影响核心功能", zap.Error(err))
	}
	a.redis = rs

	return nil
}

// initDependencies 初始化依赖注入
func (a *App) initDependencies() {
	// ========== 创建 Redis 仓储 ==========
	redisRepo := repository.NewRedisRepository(a.redis)

	// ========== 创建存储（oss / local）==========
	store, err := storage.New(a.cfg.Storage)
	if err != nil {
		logger.Fatal("存储初始化失败", zap.Error(err))
	}

	// ========== 创建 Repository ==========
	siteRepo := repository.NewSiteRepository(a.mysqlDB)
	siteAdminRepo := repository.NewSiteRepository(a.mysqlDB)
	authRepo := repository.NewAuthRepository(a.mysqlDB)
	adminAuthRepo := repository.NewAdminAuthRepository(a.mysqlDB)
	articlesRepo := repository.NewArticlesRepository(a.mysqlDB, redisRepo)
	articlesAdminRepo := repository.NewArticlesAdminRepository(a.mysqlDB)
	categoriesRepo := repository.NewCategoriesRepository(a.mysqlDB)
	categoriesAdminRepo := repository.NewCategoriesAdminRepository(a.mysqlDB)
	tagsRepo := repository.NewTagsRepository(a.mysqlDB)
	tagsAdminRepo := repository.NewTagsAdminRepository(a.mysqlDB)
	archivesRepo := repository.NewArchivesRepository(a.mysqlDB)
	statsRepo := repository.NewStatsRepository(a.mysqlDB)
	aboutRepo := repository.NewAboutRepository(a.mysqlDB)
	commentsRepo := repository.NewCommentsRepository(a.mysqlDB)
	likeRepo := repository.NewLikeRepository(a.mysqlDB)

	// ========== 创建 Service ==========
	siteService := service.NewSiteService(siteRepo)
	siteAdminService := service.NewSiteService(siteAdminRepo)
	authService := service.NewAuthService(authRepo, redisRepo)
	adminAuthService := service.NewAdminAuthService(adminAuthRepo, redisRepo)
	articlesService := service.NewArticlesService(articlesRepo)
	articlesAdminService := service.NewArticlesAdminService(articlesAdminRepo)
	categoriesService := service.NewCategoriesService(categoriesRepo)
	categoriesAdminService := service.NewCategoriesAdminService(categoriesAdminRepo)
	tagsService := service.NewTagsService(tagsRepo)
	tagsAdminService := service.NewTagsAdminService(tagsAdminRepo)
	archivesService := service.NewArchivesService(archivesRepo)
	statsService := service.NewStatsService(statsRepo)
	aboutService := service.NewAboutService(aboutRepo)
	commentsService := service.NewCommentsService(commentsRepo)
	likeService := service.NewLikeService(likeRepo)
	uploadService := service.NewUploadService(store)

	//// ========== 创建 Controller ==========
	siteController := web.NewSiteController(siteService)
	siteAdminController := admin.NewSiteAdminController(siteAdminService, redisRepo)
	adminAuthController := admin.NewAdminAuthController(adminAuthService, redisRepo)
	authController := web.NewAuthController(authService, redisRepo)
	articlesController := web.NewArticlesController(articlesService)
	adminArticlesController := admin.NewAdminArticlesController(articlesAdminService, redisRepo)
	categoriesController := web.NewCategoriesController(categoriesService)
	adminCategoriesController := admin.NewAdminCategoriesController(categoriesAdminService, redisRepo)
	tagsController := web.NewTagsController(tagsService)
	adminTagsController := admin.NewAdminTagsController(tagsAdminService, redisRepo)
	archivesController := web.NewArchivesController(archivesService)
	statsController := web.NewStatsController(statsService)
	aboutController := web.NewAboutController(aboutService)
	commentsController := web.NewCommentsController(commentsService, redisRepo, commentsRepo)
	likeController := web.NewLikeController(likeService, redisRepo)
	uploadController := admin.NewUploadController(uploadService, redisRepo)

	// ========== 创建 Router ==========
	a.router = api.NewRouter(siteController, siteAdminController,
		adminAuthController, authController, articlesController,
		adminArticlesController, categoriesController,
		adminCategoriesController, tagsController,
		adminTagsController, archivesController, statsController,
		aboutController, commentsController, likeController,
		uploadController)
}

// initRouter 初始化路由
func (a *App) initRouter() {
	// 设置 Gin 模式
	gin.SetMode(a.cfg.App.Mode)
}

// initServer 初始化 HTTP 服务器
func (a *App) initServer() {
	engine := gin.New()

	// 注册路由
	a.router.Setup(engine)

	// 创建 HTTP 服务器
	a.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler:        engine,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

// Run 运行应用
func (a *App) Run() {
	// 启动 HTTP 服务器
	go func() {
		logger.Info("HTTP 服务器启动",
			zap.String("addr", a.server.Addr),
			zap.String("mode", a.cfg.App.Mode),
		)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP 服务器启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	a.gracefulShutdown()
}

// gracefulShutdown 优雅关闭
func (a *App) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	}

	// 关闭路由连接
	if a.router != nil {
		if err := a.router.Close(); err != nil {
			logger.Error("关闭路由连接失败", zap.Error(err))
		}
	}

	// 关闭数据库连接
	_ = database.CloseMySQL()
	_ = database.CloseRedis()

	// 同步日志
	_ = logger.Sync()

	logger.Info("服务器已关闭")
	logger.Info("=========================================")
}
