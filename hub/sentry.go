package hub

import (
	"github.com/getsentry/sentry-go"
	"wechat-mp-server/config"
)

func initSentry() {
	if !enableSentry() {
		logger.Info("sentry disabled")
		return
	}
	err := sentry.Init(sentry.ClientOptions{
		Release:          "wechat-mp-server@" + config.Version,
		Dsn:              config.GlobalConfig.GetString("sentry.dsn"),
		TracesSampleRate: 0.01,
	})
	if err != nil {
		logger.Fatalf("sentry.Init: %s", err)
	}
	logger.Info("sentry init success")
	//sentry.CaptureMessage("It works!")
}

func enableSentry() bool {
	return config.GlobalConfig.InConfig("sentry.dsn")
}
