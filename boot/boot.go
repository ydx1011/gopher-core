package boot

import (
	"flag"
	"github.com/ydx1011/gopher-core"
	"github.com/ydx1011/gopher-core/bean"
	"os"
	"sync"
)

const EnvNameConfigFile = "GOPHER_CONFIG_FILE"

var (
	// 默认的配置路径
	ConfigPath = "application.yaml"

	creator func() gopher.Application = defaultCreator
	gApp    gopher.Application
	once    sync.Once
)

// 注册到全局Application
// 注册对象
// 支持注册
//  1、interface、struct指针，注册名称使用【类型名称】；
//  2、struct/interface的构造函数 func() TYPE，注册名称使用【返回值的类型名称】。
// opts添加bean注册的配置，详情查看bean.RegisterOpt
func RegisterBean(o interface{}, opts ...bean.RegisterOpt) error {
	return instance().RegisterBean(o, opts...)
}

// 注册到全局Application
// 使用指定名称注册对象
// 支持注册
//  1、interface、struct指针，注册名称使用【类型名称】；
//  2、struct/interface的构造函数 func() TYPE，注册名称使用【返回值的类型名称】。
// opts添加bean注册的配置，详情查看bean.RegisterOpt
func RegisterBeanByName(name string, o interface{}, opts ...bean.RegisterOpt) error {
	return instance().RegisterBeanByName(name, o, opts...)
}

// 自定义启动的Application
// 必须在注册对象和Run之前调用
func Customize(app gopher.Application) {
	creator = func() gopher.Application {
		return app
	}
}

func defaultCreator() gopher.Application {
	if conf, ok := os.LookupEnv(EnvNameConfigFile); ok {
		ConfigPath = conf
	}
	flag.StringVar(&ConfigPath, "f", ConfigPath, "Application configuration file path.")
	flag.Parse()
	return gopher.NewFileConfigApplication(ConfigPath)
}

func instance() gopher.Application {
	once.Do(func() {
		gApp = creator()
	})
	return gApp
}

// 启动全局Application
func Run() error {
	return instance().Run()
}
