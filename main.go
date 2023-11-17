package main

import (
	"flag"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"olympus-medusa/common"
	"olympus-medusa/config"
	routers "olympus-medusa/router"
)

var (
	tomlFile = flag.String("config", "./configs/application.yaml", "config file")
)

// Run 运行
// https://github.com/kaiyuan10nian/blog/blob/main/route/routes.go
func main() {
	// 加载配置
	InitConfig()
	// 初始化日志配置
	// 初始化数据
	_ = common.InitDB()
	// 初始化web服务
	initWeb()
}

func InitConfig() {
	config.InitLogger()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./resource/")
	err := viper.ReadInConfig()
	if err != nil {
		panic("" + err.Error())
	}
}

func initWeb() {
	gin.SetMode(gin.DebugMode)
	application := gin.Default()
	//app.NoRoute(middleware.NoRouteHandler())
	//app.NoMethod(middleware.NoMethodHandler())
	//// 崩溃恢复
	//app.Use(middleware.RecoveryMiddleware())
	// 注册路由
	routers.RegisterRouterSys(application.Group("/api"))
	//application.StaticFS("/kaiyuan", http.Dir("/opt/server/nginx-1.18/html/kaiyuan"))
	port := viper.GetString("web.port") //这里加载配置文件中的端口
	if port != "" {
		panic(application.Run(":" + port))
	}
	panic(application.Run())
}
