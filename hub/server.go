package hub

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/sirupsen/logrus"
	"sync"
	"wechat-mp-server/config"
)

var Version = "debug"

var Instance *Server

type Server struct {
	HttpEngine   *gin.Engine
	WechatEngine *officialaccount.OfficialAccount
	MsgEngine    *MsgEngine
	// 如果未来框架想要接入数据库等 可直接在此处添加
}

var logger = logrus.WithField("hub", "internal")

// Init 快速初始化
func Init() {
	logger.Info("wechat-mp-server version: ", Version)

	initSentry()

	// 初始化网络服务
	logger.Info("start init gin...")
	gin.SetMode(gin.ReleaseMode)
	httpEngine := gin.New()
	httpEngine.Use(ginRequestLog)
	if enableSentry() {
		httpEngine.Use(sentrygin.New(sentrygin.Options{}))
		logger.Info("sentry enabled")
	}

	// 初始化微信相关
	logger.Info("start init wechat...")
	wc := wechat.NewWechat()
	memoryCache := cache.NewMemory()

	cfg := &offConfig.Config{
		AppID:          config.GlobalConfig.GetString("wechat.appID"),
		AppSecret:      config.GlobalConfig.GetString("wechat.appSecret"),
		Token:          config.GlobalConfig.GetString("wechat.token"),
		EncodingAESKey: config.GlobalConfig.GetString("wechat.encodingAESKey"),
		Cache:          memoryCache,
	}
	wcOfficialAccount := wc.GetOfficialAccount(cfg)

	Instance = &Server{
		HttpEngine:   httpEngine,
		WechatEngine: wcOfficialAccount,
		MsgEngine:    NewMsgEngine(),
	}
	Instance.MsgEngine.Use(wechatMsgLog, wechatLongMsgHandle) // 注册log中间件
}

// Run 正式开启服务
func Run() {
	go func() {
		logger.Info("http engine starting...")
		if err := Instance.HttpEngine.Run(":" + config.GlobalConfig.GetString("httpEngine.port")); err != nil {
			logger.Fatal(err)
		} else {
			logger.Info("http engine running...")
		}
	}()
}

// StartService 启动服务
// 根据 Module 生命周期 此过程应在Login前调用
// 请勿重复调用
func StartService() {
	logger.Infof("initializing modules ...")
	for _, mi := range Modules {
		mi.Instance.Init()
	}
	for _, mi := range Modules {
		mi.Instance.PostInit()
	}
	logger.Info("all modules initialized")

	logger.Info("register modules serve functions ...")

	Instance.HttpEngine.Any("/serve", Instance.MsgEngine.Serve) //处理推送消息以及事件
	for _, mi := range Modules {
		mi.Instance.Serve(Instance)
	}
	logger.Info("all modules serve functions registered")

	logger.Info("starting modules tasks ...")
	for _, mi := range Modules {
		go mi.Instance.Start(Instance)
	}

	logger.Info("tasks running")
}

// Stop 停止所有服务
// 调用此函数并不会使服务器关闭
func Stop() {
	logger.Warn("stopping ...")
	wg := sync.WaitGroup{}
	for _, mi := range Modules {
		wg.Add(1)
		go mi.Instance.Stop(Instance, &wg)
	}
	wg.Wait()
	logger.Info("stopped")
	Modules = make(map[string]ModuleInfo)
}
