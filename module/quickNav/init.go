package quickNav

// 在这里做一个模块设计的简述
// 1. 为什么不把设置存入数据库？ 因为菜单中含有部分动态元素 其名称与URL等配置会动态变化，使用数据库存储此类数据复杂度较高，需要一定的低代码设计
// 2. 为什么把Item做成接口？直接设为struct不好吗？ 为了灵活的同时提高复用且能让IDE进行检查并自动补全。
// 3. 基础类型 如链接、纯文字、分割线 已在common中实现 可以直接使用 新建特殊动态按钮请添加至custom中

import (
	"sync"
	"wechat-mp-server/hub"
)

type Mod struct {
}

func (m *Mod) GetModuleInfo() hub.ModuleInfo {
	return hub.ModuleInfo{
		ID:       hub.NewModuleID("atom", "quickNav"),
		Instance: m,
	}
}

func (m *Mod) Init() {
}

func (m *Mod) PostInit() {
}

func (m *Mod) Serve(s *hub.Server) {
	s.MsgEngine.Group("快捷导航", menuHandler).MsgText("", 3).EventClick("")
}

func (m *Mod) Start(_ *hub.Server) {
}

func (m *Mod) Stop(_ *hub.Server, wg *sync.WaitGroup) {
	wg.Done()
}
