package main

import (
	"os"
	"os/signal"
	"syscall"
	"wechat-mp-server/config"
	"wechat-mp-server/hub" // 框架
	"wechat-mp-server/module/steggoDemo"
	"wechat-mp-server/utils"

	"wechat-mp-server/module/pong"
	"wechat-mp-server/module/quickNav"
	"wechat-mp-server/module/templateMessage"
	"wechat-mp-server/module/timeoutTest"
	"wechat-mp-server/module/wechatPong"
	//_ "wechat-mp-server/module/wechatApiProxy"  // 微信管理员代理模块 用于代理Wechat相关api接口并自动加上accessToken
)

func init() {
	utils.WriteLogToFS()
	config.Init()
}

func main() {
	// 新增module后请在下方引入
	hub.RegisterModule(
		&pong.Mod{},            // gin的ping-pong模块
		&wechatPong.Mod{},      // 微信ping-pong模块
		&quickNav.Mod{},        // 快捷导航模块 用于解决微信菜单按钮不足的情况
		&templateMessage.Mod{}, // 提供发送模板消息api
		&timeoutTest.Mod{},     // 超时回复测试模块
		&steggoDemo.Mod{},      // steggo的应用案例
	)

	hub.Init()

	hub.StartService()

	hub.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	hub.Stop()
}
