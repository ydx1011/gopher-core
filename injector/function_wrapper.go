package injector

import (
	"errors"
	"fmt"
	"github.com/ydx1011/gopher-core/bean"
	"reflect"
)

func WrapBeanFactoryByNameFunc(o interface{}, names []string, container bean.Container, injector Injector) (interface{}, error) {
	ft := reflect.TypeOf(o)
	if ft.Kind() != reflect.Func {
		return o, nil
	}
	if ft.NumOut() != 1 {
		return o, fmt.Errorf("Bean Factory function: %s without any return value ", ft.String())
	}

	rt := ft.Out(0)
	if rt.Kind() != reflect.Ptr && rt.Kind() != reflect.Interface {
		return o, fmt.Errorf("Bean Factory function: %s 1st return value must be pointer or interface ", ft.String())
	}
	pn := ft.NumIn()
	if pn > 0 {
		if pn != len(names) {
			return o, fmt.Errorf("Bean Factory function: %s have %d params but with %d names, Not match ", ft.String(), pn, len(names))
		}
		retFv := reflect.MakeFunc(reflect.FuncOf(nil, []reflect.Type{rt}, false), func(args []reflect.Value) (results []reflect.Value) {
			fv := reflect.ValueOf(o)
			values := make([]reflect.Value, pn)
			for i := 0; i < pn; i++ {
				o := reflect.New(ft.In(i)).Elem()
				err := injector.InjectValue(container, names[i], o)
				if err != nil {
					err = fmt.Errorf("Inject function [%s] param %d [%s] failed:error: %s\n", ft.String(), i, o.Type().String(), err.Error())
					panic(err)
				}
				values[i] = o
			}

			return fv.Call(values)
		})
		return retFv.Interface(), nil
	}
	return o, nil
}

func WrapBeanFactoryFunc(o interface{}, container bean.Container, injector Injector) (interface{}, error) {
	ft := reflect.TypeOf(o)
	if ft.Kind() != reflect.Func {
		return o, nil
	}
	if ft.NumOut() != 1 {
		return o, fmt.Errorf("Bean Factory function: %s without any return value ", ft.String())
	}

	rt := ft.Out(0)
	if rt.Kind() != reflect.Ptr && rt.Kind() != reflect.Interface {
		return o, errors.New("Bean function 1st return value must be pointer or interface. ")
	}
	pn := ft.NumIn()
	if pn > 0 {
		retFv := reflect.MakeFunc(reflect.FuncOf(nil, []reflect.Type{rt}, false), func(args []reflect.Value) (results []reflect.Value) {
			fv := reflect.ValueOf(o)
			values := make([]reflect.Value, pn)
			for i := 0; i < pn; i++ {
				o := reflect.New(ft.In(i)).Elem()
				err := injector.InjectValue(container, "", o)
				if err != nil {
					err = fmt.Errorf("Inject function [%s] failed:error: %s\n", ft.Name(), err.Error())
					panic(err)
				}
				values[i] = o
			}

			return fv.Call(values)
		})
		return retFv.Interface(), nil
	}
	return o, nil
}
