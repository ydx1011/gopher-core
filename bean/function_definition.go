package bean

import (
	"errors"
	"reflect"
)

func verifyBeanFunctionEx(ft reflect.Type) error {
	if ft.Kind() != reflect.Func {
		return errors.New("Param not function ")
	}
	if ft.NumOut() != 1 {
		return errors.New("Bean function must ONLY have 1 return value ")
	}

	rt := ft.Out(0)
	if rt.Kind() != reflect.Ptr && rt.Kind() != reflect.Interface {
		return errors.New("Bean function 1st return value must be pointer or interface ")
	}

	return nil
}
