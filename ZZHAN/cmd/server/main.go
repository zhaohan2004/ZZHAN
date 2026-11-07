package main

import (
	"ZZHAN/internal/app"
	"flag"
	"fmt"
	"log"
)

func main() {
	// 定义命令行参数
	configPath := flag.String("c", "", "配置文件路径")
	migrateCovers := flag.Bool("migrate-covers", false, "将历史内嵌封面迁移到配置的图片存储")
	flag.Parse()

	if *migrateCovers {
		if err := runCoverMigration(*configPath); err != nil {
			log.Fatal(err)
		}
		return
	}

	// 创建应用实例
	application := app.NewApp()

	// 初始化应用
	if err := application.Initialize(*configPath); err != nil {
		panic(fmt.Sprintf("应用初始化失败: %v", err))
	}

	// 运行应用
	application.Run()
}
