package bean

type Initializing interface {
	// 当初始化和注入完成时回调
	BeanAfterSet() error
}

type Disposable interface {
	// 进入销毁阶段，应该尽快做回收处理并退出处理任务
	BeanDestroy() error
}
