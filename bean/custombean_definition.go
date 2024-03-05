package bean

import (
	"fmt"
	"github.com/ydx1011/gopher-core/reflection"
	"reflect"
)

type CustomBeanFactory interface {
	// 返回或者创建bean的方法
	// 该方法可能包含一个或者多个参数，参数会在实例化时自动注入
	// 该方法只能有一个返回值，返回的值将被注入到依赖该类型值的对象中
	BeanFactory() interface{}

	// BeanFactory返回创建bean方法如果带参数，且参数需要指定注入名称时将根据InjectNames返回的名称列表进行匹配
	// 注意：
	// 1、如果所有参数都不需要名称匹配，则返回nil
	// 2、如果需要使用名称匹配则：返回的string数组长度需要与创建bean方法的常数个数一致
	// 3、如果需要部分匹配，则需要自动匹配的参数对应的name填入空字符串""
	InjectNames() []string

	// BeanFactory返回值包含的初始化方法名，可为空
	InitMethodName() string

	// BeanFactory返回值包含的销毁方法名，可为空
	DestroyMethodName() string
}

type defaultCustomBeanFactory struct {
	beanFunc      interface{}
	names         []string
	initMethod    string
	destroyMethod string
}

func NewCustomBeanFactory(beanFunc interface{}, initMethod, destroyMethod string) *defaultCustomBeanFactory {
	ft := reflect.TypeOf(beanFunc)
	if err := verifyBeanFunctionEx(ft); err != nil {
		panic(fmt.Errorf("NewCustomMethodBean with a invalid function type: %s, error: %v", ft.String(), err))
	}
	return &defaultCustomBeanFactory{
		beanFunc:      beanFunc,
		initMethod:    initMethod,
		destroyMethod: destroyMethod,
	}
}

func NewCustomBeanFactoryWithName(beanFunc interface{}, names []string, initMethod, destroyMethod string) *defaultCustomBeanFactory {
	ft := reflect.TypeOf(beanFunc)
	if err := verifyBeanFunctionEx(ft); err != nil {
		panic(fmt.Errorf("NewCustomMethodBean with a invalid function type: %s", ft.String()))
	}
	return &defaultCustomBeanFactory{
		beanFunc:      beanFunc,
		names:         names,
		initMethod:    initMethod,
		destroyMethod: destroyMethod,
	}
}

func (b *defaultCustomBeanFactory) BeanFactory() interface{} {
	return b.beanFunc
}

func (b *defaultCustomBeanFactory) InjectNames() []string {
	return b.names
}

func (b *defaultCustomBeanFactory) InitMethodName() string {
	return b.initMethod
}

func (b *defaultCustomBeanFactory) DestroyMethodName() string {
	return b.destroyMethod
}

type customMethodBeanDefinition struct {
	functionExDefinition
	initializingFuncName string
	disposableFuncName   string
}

func newCustomMethodBeanDefinition(b CustomBeanFactory) (Definition, error) {
	d, err := newFunctionExDefinition(b.BeanFactory())
	if err != nil {
		return nil, err
	}
	ret := &customMethodBeanDefinition{
		functionExDefinition: *d.(*functionExDefinition),
		initializingFuncName: b.InitMethodName(),
		disposableFuncName:   b.DestroyMethodName(),
	}

	return ret, ret.verifyCustomBeanFunction()
}

func (d *customMethodBeanDefinition) verifyCustomBeanFunction() error {
	rt := d.t
	if d.initializingFuncName != "" {
		if !checkPublic(d.initializingFuncName) {
			return fmt.Errorf("Type %s init method %s is private ", reflection.GetTypeName(d.t), d.initializingFuncName)
		}
		m, ok := rt.MethodByName(d.initializingFuncName)
		if !ok {
			return fmt.Errorf("Type %s init method %s not found ", reflection.GetTypeName(d.t), d.initializingFuncName)
		} else {
			if m.Type.NumIn() == 0 {
				return fmt.Errorf("Type %s init method %s cannot with params ", reflection.GetTypeName(d.t), d.initializingFuncName)
			}
		}
	}

	if d.disposableFuncName != "" {
		if !checkPublic(d.initializingFuncName) {
			return fmt.Errorf("Type %s destroy method %s is private ", reflection.GetTypeName(d.t), d.initializingFuncName)
		}
		m, ok := rt.MethodByName(d.disposableFuncName)
		if !ok {
			return fmt.Errorf("Type %s destroy method %s not found ", reflection.GetTypeName(d.t), d.disposableFuncName)
		} else {
			if m.Type.NumIn() == 0 {
				return fmt.Errorf("Type %s destroy method %s cannot with params ", reflection.GetTypeName(d.t), d.disposableFuncName)
			}
		}
	}

	return nil
}

func checkPublic(name string) bool {
	return name[0] >= 'A' && name[0] <= 'Z'
}
