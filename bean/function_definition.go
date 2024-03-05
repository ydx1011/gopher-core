package bean

import (
	"errors"
	"fmt"
	errors2 "github.com/ydx1011/gopher-core/errors"
	"github.com/ydx1011/gopher-core/reflection"
	"reflect"
	"sync"
	"sync/atomic"
)

const (
	functionDefinitionNone = iota
	functionDefinitionInjecting
)

var (
	DummyType  = reflect.TypeOf((*struct{})(nil)).Elem()
	DummyValue = reflect.ValueOf(struct{}{})
)

type functionExDefinition struct {
	name   string
	o      interface{}
	fn     reflect.Value
	t      reflect.Type
	status int32

	instances    reflect.Value
	instanceLock sync.RWMutex
	initOnce     int32
	destroyOnce  int32
}

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

func newFunctionExDefinition(o interface{}) (Definition, error) {
	ft := reflect.TypeOf(o)
	err := verifyBeanFunctionEx(ft)
	if err != nil {
		return nil, err
	}
	ot := ft.Out(0)
	fn := reflect.ValueOf(o)
	ret := &functionExDefinition{
		o:         o,
		name:      reflection.GetTypeName(ot),
		fn:        fn,
		t:         ot,
		instances: reflect.MakeMap(reflect.MapOf(ot, DummyType)),
	}
	return ret, nil
}

func (d *functionExDefinition) Type() reflect.Type {
	return d.t
}

func (d *functionExDefinition) Name() string {
	return d.name
}

func (d *functionExDefinition) Value() reflect.Value {
	if atomic.CompareAndSwapInt32(&d.status, functionDefinitionNone, functionDefinitionInjecting) {
		defer atomic.CompareAndSwapInt32(&d.status, functionDefinitionInjecting, functionDefinitionNone)
		v := d.fn.Call(nil)[0]
		if v.IsValid() {
			d.instanceLock.Lock()
			defer d.instanceLock.Unlock()
			d.instances.SetMapIndex(v, DummyValue)
		}
		return v
	} else {
		panic(fmt.Errorf("BeanDefinition: [Function] inject type [%s] Circular dependency ", d.name))
	}
}

func (d *functionExDefinition) Interface() interface{} {
	return d.o
}

func (d *functionExDefinition) IsObject() bool {
	return false
}

func (d *functionExDefinition) AfterSet() error {
	if atomic.CompareAndSwapInt32(&d.initOnce, 0, 1) {
		d.instanceLock.RLock()
		defer d.instanceLock.RUnlock()
		var errs errors2.Errors

		for _, i := range d.instances.MapKeys() {
			if i.IsValid() && !i.IsNil() {
				if v, ok := i.Interface().(Initializing); ok {
					err := v.BeanAfterSet()
					if err != nil {
						_ = errs.AddError(err)
					}
				}
			}
		}
		if errs.Empty() {
			return nil
		}
		return errs
	}
	return nil
}

func (d *functionExDefinition) Destroy() error {
	if atomic.CompareAndSwapInt32(&d.destroyOnce, 0, 1) {
		d.instanceLock.RLock()
		defer d.instanceLock.RUnlock()
		var errs errors2.Errors
		for _, i := range d.instances.MapKeys() {
			if i.IsValid() && !i.IsNil() {
				if v, ok := i.Interface().(Disposable); ok {
					err := v.BeanDestroy()
					if err != nil {
						_ = errs.AddError(err)
					}
				}
			}
		}
		if errs.Empty() {
			return nil
		}
		return errs
	}
	return nil
}

func (d *functionExDefinition) Classify(classifier Classifier) (bool, error) {
	d.instanceLock.RLock()
	defer d.instanceLock.RUnlock()
	var errs errors2.Errors
	ok := false
	for _, i := range d.instances.MapKeys() {
		if !i.IsNil() {
			ret, err := classifier.Classify(i.Interface())
			if ret {
				ok = ret
			}
			if err != nil {
				_ = errs.AddError(err)
			}
		}
	}
	if errs.Empty() {
		return ok, nil
	}
	return ok, errs
}
