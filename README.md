# wechat-mp-server
Go版本的微信公众号后端框架 包含消息路由与基础的消息处理中间件

是一个多模块组合设计的微信公众号后端程序

包装了基础功能,同时设计了一个~~良好~~的项目结构

## Demo

搜索微信公众号：`宅男的天台`

支持指令：
> ping


## 不了解go?

golang 极速入门

[点我看书](https://github.com/justjavac/free-programming-books-zh_CN#go)

## Module 配置

module参考[wechatPong](https://github.com/asjdf/wechat-mp-server/blob/main/module/wechatPong/init.go)

编写自己的Module后在[main.go](https://github.com/asjdf/wechat-mp-server/blob/main/main.go)中启用Module

## 功能
- [x] 消息路由
- [x] 长消息分段发送
- [x] 超时回复处理
- [x] 简单的模组定义

## 快速开始
1. 根据 application.yaml.example 的结构编写 application.yaml，框架默认的被动回复消息接口在 /serve
2. 根据 application 的配置创建反代或内网穿透
3. 构建并运行
4. 向你的微信公众号发送“ping”、“快捷导航”、“超时回复测试”、“超时回复测试2”，即可测试框架基本功能

