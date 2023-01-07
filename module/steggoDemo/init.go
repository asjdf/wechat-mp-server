package steggoDemo

import (
	"sync"
	"wechat-mp-server/hub"
)

type Mod struct {
}

func (m *Mod) GetModuleInfo() hub.ModuleInfo {
	//TODO implement me
	panic("implement me")
}

func (m *Mod) Init() {
	//TODO implement me
	panic("implement me")
}

func (m *Mod) PostInit() {
	//TODO implement me
	panic("implement me")
}

func (m *Mod) Serve(s *hub.Server) {
	//TODO implement me
	panic("implement me")
}

func (m *Mod) Start(s *hub.Server) {
	//TODO implement me
	panic("implement me")
}

func (m *Mod) Stop(s *hub.Server, wg *sync.WaitGroup) {
	wg.Done()
}

