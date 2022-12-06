package hub

import (
	"context"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"regexp"
	"sync"
)

// Message 最初只是为了加两个方法上去 现在作为上下文
type Message struct {
	*message.MixMessage
	Reply *message.Reply

	// Keys 是每个请求所特有的键值对
	Keys map[string]interface{}
	// 这个锁保护 Keys map
	sync.RWMutex

	//记录当前handler序号
	Index int8

	handlers HandlersChain

	//用于传递sentry
	Span context.Context

	//用于记录触发的key
	Pattern *regexp.Regexp
}

func (m *Message) Type() (msgType string) {
	if m.MsgType == message.MsgTypeEvent {
		msgType = string(m.Event)
	} else {
		msgType = string(m.MsgType)
	}
	return
}

func (m *Message) Key() (key string) {
	if m.MsgType == message.MsgTypeEvent {
		key = m.EventKey
	} else {
		key = m.Content
	}
	return
}

func (m *Message) Next() {
	m.Index++
	for m.Index < int8(len(m.handlers)) {
		m.handlers[m.Index](m)
		m.Index++
	}
}

func (m *Message) Abort() {
	m.Index = int8(len(m.handlers))
}

func (m *Message) Set(key string, value interface{}) {
	m.Lock()
	if m.Keys == nil {
		m.Keys = make(map[string]interface{})
	}

	m.Keys[key] = value
	m.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
func (m *Message) Get(key string) (value interface{}, exists bool) {
	m.RLock()
	value, exists = m.Keys[key]
	m.RUnlock()
	return
}

// MustGet 如果键存在，则返回给定键的值，否则panic
func (m *Message) MustGet(key string) interface{} {
	if value, exists := m.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString 返回string类型的指定键的值
func (m *Message) GetString(key string) (s string) {
	if val, ok := m.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}
