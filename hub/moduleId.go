package hub

import "strings"

// ModuleID 模块ID
// 请使用 小写 并用 _ 代替空格
// Example:
// - atom.pong
type ModuleID interface {
	Namespace() string  // 命名空间，一般用开发者名字，请使用 小写 并用 _ 代替空格
	ModuleName() string // 模块名，一般用包名
	String() string
}

type moduleIDStruct struct {
	namespace  string // 命名空间，一般用开发者名字，请使用 小写 并用 _ 代替空格
	moduleName string // 模块名，一般用包名
}

func (id moduleIDStruct) Namespace() string {
	return id.namespace
}

func (id moduleIDStruct) ModuleName() string {
	return id.moduleName
}

// String 实现 Stringer 接口
// 凡是将 ModuleID 转成 string 一律用这个接口, 无论是输出还是获取模块名
func (id moduleIDStruct) String() string {
	return id.namespace + "." + id.moduleName
}

// NewModuleID 构造函数，统一生成 ModuleID，避免非法模块名的存在
func NewModuleID(namespace, moduleName string) ModuleID {
	namespace = strings.TrimSpace(namespace)
	moduleName = strings.TrimSpace(moduleName)
	if namespace == "" || moduleName == "" {
		panic("模块名的namespace和name均不能为空白")
	}
	return moduleIDStruct{namespace: namespace, moduleName: moduleName}
}
