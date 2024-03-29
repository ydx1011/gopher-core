package gopher

import (
	"github.com/xfali/xlog"
	"github.com/ydx1011/gopher-core/appcontext"
	"github.com/ydx1011/gopher-core/bean"
	"github.com/ydx1011/gopher-core/util"
	"github.com/ydx1011/yfig"
)

type Application interface {
	// 注册对象
	// 支持注册
	//  1、interface、struct指针，注册名称使用【类型名称】；
	//  2、struct/interface的构造函数 func() TYPE，注册名称使用【返回值的类型名称】。
	// opts添加bean注册的配置，详情查看bean.RegisterOpt
	RegisterBean(o interface{}, opts ...RegisterOpt) error

	// 使用指定名称注册对象
	// 支持注册
	//  1、interface、struct指针，注册名称使用【类型名称】；
	//  2、struct/interface的构造函数 func() TYPE，注册名称使用【返回值的类型名称】。
	// opts添加bean注册的配置，详情查看bean.RegisterOpt
	RegisterBeanByName(name string, o interface{}, opts ...RegisterOpt) error

	AddListeners(listeners ...interface{})

	// 启动应用容器
	Run() error
}

type RegisterOpt = bean.RegisterOpt

type FileConfigApplication struct {
	ctx    appcontext.ApplicationContext
	logger xlog.Logger
}

type Opt func(*FileConfigApplication)

func NewFileConfigApplication(configPath string, opts ...Opt) *FileConfigApplication {
	// Disable fig's log
	//yfig.SetLog(func(format string, o ...interface{}) {})
	prop, err := yfig.LoadYamlFile(configPath)
	if err != nil {
		xlog.Errorln("load config file failed: ", err)
		return nil
	}
	return NewApplication(prop, opts...)
}

func NewApplication(prop yfig.Properties, opts ...Opt) *FileConfigApplication {
	if prop == nil {
		xlog.Errorln("Properties cannot be nil. ")
		return nil
	}
	ret := &FileConfigApplication{
		ctx:    appcontext.NewDefaultApplicationContext(),
		logger: xlog.GetLogger(),
	}

	for _, opt := range opts {
		opt(ret)
	}

	err := ret.ctx.Init(prop)
	if err != nil {
		ret.logger.Fatalln(err)
		return nil
	}

	return ret
}

func (app *FileConfigApplication) RegisterBean(o interface{}, opts ...RegisterOpt) error {
	return app.ctx.RegisterBean(o, opts...)
}

func (app *FileConfigApplication) RegisterBeanByName(name string, o interface{}, opts ...RegisterOpt) error {
	return app.ctx.RegisterBeanByName(name, o, opts...)
}

func (app *FileConfigApplication) AddListeners(listeners ...interface{}) {
	app.ctx.AddListeners(listeners...)
}

func (app *FileConfigApplication) Run() error {
	err := app.ctx.Start()
	if err != nil {
		return err
	}
	return util.HandlerSignal(app.logger, app.ctx.Close)
}
