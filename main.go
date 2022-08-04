package main

import (
	"os"
	"os/signal"
	"syscall"
	"wechat-mp-server/config"
	"wechat-mp-server/hub" // 框架
	"wechat-mp-server/utils"

	// 新增module后请在下方引入
	_ "wechat-mp-server/module/pong"            // gin的ping-pong模块
	_ "wechat-mp-server/module/quickNav"        // 快捷导航模块 用于解决微信菜单按钮不足的情况
	_ "wechat-mp-server/module/templateMessage" // 提供发送模板消息api
	_ "wechat-mp-server/module/timeoutTest"     // 超时回复测试模块
	//_ "wechat-mp-server/module/wechatApiProxy"  // 微信管理员代理模块 用于代理Wechat相关api接口并自动加上accessToken
	_ "wechat-mp-server/module/wechatPong" // 微信ping-pong模块
)

func init() {
	utils.WriteLogToFS()
	config.Init()
}

func main() {
	hub.Init()

	hub.StartService()

	hub.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	hub.Stop()
}
