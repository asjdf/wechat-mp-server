package hub

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"reflect"
	"regexp"
	"sort"
	"sync"
)

type HandlerFunc func(msg *Message)

type HandlersChain []HandlerFunc

// IRoutes 定义所有的路由处理接口
type IRoutes interface {
	Group(baseKey string, middleware ...HandlerFunc) *RouterGroup
	Use(middleware ...HandlerFunc) IRoutes

	Handle(handle interface{})

	// 下面方法的命名规则按照事件大类+事件小类的结构，方便自动补全时能快速找到已实现的相关功能

	MsgText(key string, index int, handler ...HandlerFunc) IRoutes
	EventClick(key string, handler ...HandlerFunc) IRoutes
	EventView(key string, handler ...HandlerFunc) IRoutes
	EventScan(index int, handler ...HandlerFunc) IRoutes
	EventSubscribe(index int, handler ...HandlerFunc) IRoutes
	EventUnsubscribe(index int, handler ...HandlerFunc) IRoutes
}

type RouterGroup struct {
	Handlers HandlersChain
	baseKey  string
	engine   *MsgEngine
	root     bool
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateKey(key string) string {
	return group.baseKey + key
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.engine
	}
	return group
}

func (group *RouterGroup) Group(baseKey string, middleware ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(middleware),
		baseKey:  group.calculateKey(baseKey),
		engine:   group.engine,
	}
}

func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

func (group *RouterGroup) Handle(h interface{}) {
	var err error
	group.engine.mu.Lock()
	switch h := h.(type) {
	case TextMsgHandle:
		group.engine.onTextMsg = append(group.engine.onTextMsg, h)
		err = sortRouter(&group.engine.onTextMsg)
	case ClickEventHandle:
		group.engine.onClickEvent = append(group.engine.onClickEvent, h)
		err = sortRouter(&group.engine.onClickEvent)
	case ViewEventHandle:
		group.engine.onViewEvent = append(group.engine.onViewEvent, h)
		err = sortRouter(&group.engine.onViewEvent)
	case ScanEventHandle:
		group.engine.onScanEvent = append(group.engine.onScanEvent, h)
		err = sortRouter(&group.engine.onScanEvent)
	case SubscribeEventHandle:
		group.engine.onSubscribeEvent = append(group.engine.onSubscribeEvent, h)
		err = sortRouter(&group.engine.onSubscribeEvent)
	case UnsubscribeEventHandle:
		group.engine.onUnsubscribeEvent = append(group.engine.onUnsubscribeEvent, h)
		err = sortRouter(&group.engine.onUnsubscribeEvent)
	default:
		logger.Errorf("Didn't support this handle: %v in %v", reflect.TypeOf(h).Name(), reflect.TypeOf(h).PkgPath())
	}
	group.engine.mu.Unlock()
	if err != nil {
		logger.Errorf("Failed to sort router: %v", err.Error())
	}
}

func (group *RouterGroup) MsgText(key string, index int, handler ...HandlerFunc) IRoutes {
	h := TextMsgHandle{}
	h.Keyword = group.calculateKey(key)
	h.Pattern = regexp.MustCompile(h.Keyword)
	h.Index = index
	h.Handlers = group.combineHandlers(handler)
	group.Handle(h)
	return group
}

func (group *RouterGroup) EventClick(key string, handler ...HandlerFunc) IRoutes {
	h := ClickEventHandle{}
	h.Keyword = group.calculateKey(key)
	h.Pattern = regexp.MustCompile(h.Keyword)
	h.Handlers = group.combineHandlers(handler)
	group.Handle(h)
	return group
}

func (group *RouterGroup) EventView(key string, handler ...HandlerFunc) IRoutes {
	h := ViewEventHandle{}
	h.Keyword = group.calculateKey(key)
	h.Pattern = regexp.MustCompile(h.Keyword)
	h.Handlers = group.combineHandlers(handler)
	group.Handle(h)
	return group
}

func (group *RouterGroup) EventScan(index int, handler ...HandlerFunc) IRoutes {
	h := ScanEventHandle{}
	h.Index = index
	h.Handlers = group.combineHandlers(handler)
	group.Handle(h)
	return group
}

func (group *RouterGroup) EventSubscribe(index int, handler ...HandlerFunc) IRoutes {
	h := SubscribeEventHandle{}
	h.Index = index
	h.Handlers = group.combineHandlers(handler)
	group.Handle(h)
	return group
}

func (group *RouterGroup) EventUnsubscribe(index int, handler ...HandlerFunc) IRoutes {
	h := UnsubscribeEventHandle{}
	h.Index = index
	h.Handlers = group.combineHandlers(handler)
	group.Handle(h)
	return group
}

var _ IRoutes = &RouterGroup{}

// BasicHandle 包含消息路由必备的内容
type BasicHandle struct {
	Index    int // 回复优先级 数字越大优先级越高
	Handlers HandlersChain
}

// BasicKeyHandle 包含消息路由必备的内容和Key
type BasicKeyHandle struct {
	Keyword string         // 匹配关键字
	Pattern *regexp.Regexp // 将 Keyword 预处理为正则对象
	BasicHandle
}

type TextMsgHandle BasicKeyHandle
type ClickEventHandle BasicKeyHandle // 本来设计是使用map来匹配这种单一结果的路由 但是为了方便使用group完成花活 所以这么搞了
type ViewEventHandle BasicKeyHandle
type ScanEventHandle BasicHandle
type SubscribeEventHandle BasicHandle
type UnsubscribeEventHandle BasicHandle

type sortItem struct {
	data interface{}
	int64
}
type sortList []sortItem

func (s sortList) Len() int {
	return len(s)
}

func (s sortList) Less(i, j int) bool {
	return s[i].int64 > s[j].int64
}

func (s sortList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func GetIndex(in interface{}) (int64, error) {
	f, ok := reflect.TypeOf(in).FieldByName("Index")
	if !ok {
		return 0, errors.New("index not found")
	}
	if f.Type.Kind() != reflect.Int {
		return 0, errors.New("index has not valid type")
	}
	return reflect.ValueOf(in).FieldByName("Index").Int(), nil
}

func iToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		panic("to slice arr not slice")
	}
	l := v.Elem().Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Elem().Index(i).Interface()
	}
	return ret
}

func sortRouter(in interface{}) error {
	rv := reflect.ValueOf(in)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New(reflect.TypeOf(in).Name() + "not ptr")
	}
	ISlice := iToSlice(in)
	structList := make(sortList, 0)
	for _, v := range ISlice {
		if index, err := GetIndex(v); err != nil {
			return err
		} else {
			structList = append(structList, sortItem{v, index})
		}
	}
	sort.Sort(structList)
	for i, v := range structList {
		rv.Elem().Index(i).Set(reflect.ValueOf(v.data))
	}
	return nil
}

type MsgEngine struct {
	RouterGroup

	onTextMsg          []TextMsgHandle
	onClickEvent       []ClickEventHandle
	onViewEvent        []ViewEventHandle
	onScanEvent        []ScanEventHandle
	onSubscribeEvent   []SubscribeEventHandle
	onUnsubscribeEvent []UnsubscribeEventHandle

	mu sync.RWMutex
}

var _ IRoutes = &MsgEngine{}

// NewMsgEngine 创建一个 MsgEngine 对象
func NewMsgEngine() *MsgEngine {
	engine := &MsgEngine{
		RouterGroup: RouterGroup{
			root: true,
		},
	}
	engine.RouterGroup.engine = engine
	return engine
}

// Serve 服务代码 服务类消息 [?] 此处应为 中间件内部调用
func (msgEngine *MsgEngine) Serve(c *gin.Context) {
	messageServer := Instance.WechatEngine.GetServer(c.Request, c.Writer)
	//messageServer.SkipValidate(true) // hack fix
	messageServer.SetMessageHandler(msgEngine.genMsgHandler())
	//处理消息接收以及回复
	err := messageServer.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	_ = messageServer.Send()
}

func (msgEngine *MsgEngine) genMsgHandler() func(msg *message.MixMessage) *message.Reply {
	return func(msg *message.MixMessage) *message.Reply {
		// 新建上下文
		superMessage := &Message{MixMessage: msg, RWMutex: sync.RWMutex{}, Index: -1}

		// 下面的内容还可简化但还未作处理 等一个有缘人来冲
		switch superMessage.Type() {
		case string(message.MsgTypeText):
			for _, h := range msgEngine.onTextMsg {
				key := superMessage.Key()
				if h.Pattern.MatchString(key) {
					superMessage.handlers = h.Handlers
					superMessage.Pattern = h.Pattern
					superMessage.Next()
				}
			}
		case string(message.EventClick):
			for _, h := range msgEngine.onClickEvent {
				key := superMessage.Key()
				if h.Pattern.MatchString(key) {
					superMessage.handlers = h.Handlers
					superMessage.Pattern = h.Pattern
					superMessage.Next()
				}
			}
		case string(message.EventView):
			for _, h := range msgEngine.onViewEvent {
				key := superMessage.Key()
				if h.Pattern.MatchString(key) {
					superMessage.handlers = h.Handlers
					superMessage.Pattern = h.Pattern
					superMessage.Next()
				}
			}
		case string(message.EventScan):
			for _, h := range msgEngine.onScanEvent {
				superMessage.handlers = h.Handlers
				superMessage.Next()
				if superMessage.Reply != nil {
					break
				}
			}
		case string(message.EventSubscribe):
			for _, h := range msgEngine.onSubscribeEvent {
				superMessage.handlers = h.Handlers
				superMessage.Next()
				if superMessage.Reply != nil {
					break
				}
			}
		case string(message.EventUnsubscribe):
			for _, h := range msgEngine.onUnsubscribeEvent {
				superMessage.handlers = h.Handlers
				superMessage.Next()
				if superMessage.Reply != nil {
					break
				}
			}
		}

		return superMessage.Reply
	}
}
