// Package wechatApiProxy 用于代理Wechat相关api接口
// 仅代理https://api.weixin.qq.com/cgi-bin/前缀的api
// 前端调用时仅需传前缀以外的路由即可 方法和微信api文档保持一致
// 如：https://api.weixin.qq.com/cgi-bin/media/upload 请求 /wechat/proxy/media/upload
package wechatApiProxy

import (
	"sync"
	"wechat-mp-server/hub"
)

func init() {
	instance = &mod{}
	hub.RegisterModule(instance)
}

var instance *mod

type mod struct {
}

func (m *mod) GetModuleInfo() hub.ModuleInfo {
	return hub.ModuleInfo{
		ID:       hub.NewModuleID("atom", "wechatApiProxy"),
		Instance: instance,
	}
}

func (m *mod) Init() {
}

func (m *mod) PostInit() {

}

func (m *mod) Serve(s *hub.Server) {
	s.HttpEngine.Any("/wechat/proxy/*route", ProxyHandler(s))
}

func (m *mod) Start(_ *hub.Server) {

}

func (m *mod) Stop(_ *hub.Server, wg *sync.WaitGroup) {
	wg.Done()
}
