package injector

import (
	"errors"
	"fmt"
	"github.com/xfali/xlog"
	"github.com/ydx1011/gopher-core/bean"
	"github.com/ydx1011/gopher-core/reflection"
	"github.com/ydx1011/gopher-core/util"
	"reflect"
	"strings"
	"sync"
)

const (
	defaultInjectTagName    = "inject"
	defaultRequiredTagField = "required"
	defaultOmitTagField     = "omiterror"
)

var (
	InjectTagName    = defaultInjectTagName
	RequiredTagField = defaultRequiredTagField
	OmitTagField     = defaultOmitTagField
)

type defaultInjector struct {
	logger    xlog.Logger
	actuators map[reflect.Kind]Actuator
	lm        ListenerManager
	tagName   string
	recursive bool
}

type Opt func(*defaultInjector)

func New(opts ...Opt) *defaultInjector {
	ret := &defaultInjector{
		logger:  xlog.GetLogger(),
		tagName: InjectTagName,
	}
	ret.actuators = map[reflect.Kind]Actuator{
		reflect.Interface: ret.injectInterface,
		reflect.Struct:    ret.injectStruct,
		reflect.Slice:     ret.injectSlice,
		reflect.Map:       ret.injectMap,
	}
	ret.lm = NewListenerManager(ret.logger)
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func (injector *defaultInjector) CanInject(o interface{}) bool {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Interface {
		return true
	}
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Struct {
			return true
		}
	}
	return false
}

func (injector *defaultInjector) Inject(c bean.Container, o interface{}) error {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Interface {
		return injector.injectInterface(c, "", v)
	}
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() == reflect.Struct {
		return injector.injectStructFields(c, v)
	}
	return errors.New("Type Not support. ")
}
func (injector *defaultInjector) CanInjectType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	actuate := injector.actuators[t.Kind()]
	return actuate != nil
}

func (injector *defaultInjector) InjectValue(c bean.Container, name string, v reflect.Value) error {
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v.CanSet() {
		actuate := injector.actuators[t.Kind()]
		if actuate == nil {
			return errors.New("Cannot inject this kind: " + t.Name())
		}
		return actuate(c, name, v)
	} else {
		return errors.New("Inject Failed: Value cannot set. ")
	}
}

type defaultListenerManager struct {
	listeners sync.Map
}

func NewListenerManager(param ...xlog.Logger) *defaultListenerManager {
	var logger xlog.Logger
	if len(param) > 0 {
		logger = param[0]
	} else {
		logger = xlog.GetLogger()
	}
	ret := &defaultListenerManager{}
	ret.listeners.Store(RequiredTagField, NewRequiredListener())
	ret.listeners.Store(OmitTagField, NewOmitErrorListener(logger))
	return ret
}

func (mgr *defaultListenerManager) AddListener(name string, listener Listener) {
	if listener != nil {
		mgr.listeners.Store(name, listener)
	}
}

func (mgr *defaultListenerManager) ParseListener(tag string) (string, []Listener) {
	strs := strings.Split(tag, ",")
	opts := strs[1:]
	// default must be required
	if len(opts) == 0 {
		opts = []string{RequiredTagField}
	}

	ret := make([]Listener, 0, len(opts))
	for _, v := range opts {
		l, ok := mgr.listeners.Load(v)
		if ok && l != nil {
			ret = append(ret, l.(Listener))
		}
	}

	return strs[0], ret
}

type RequiredListener struct{}

func NewRequiredListener() *RequiredListener {
	l := RequiredListener{}
	return &l
}

func (l *RequiredListener) OnInjectFailed(err error) {
	panic(err)
}

type OmitErrorListener struct {
	logger xlog.Logger
}

func NewOmitErrorListener(logger xlog.Logger) *OmitErrorListener {
	l := OmitErrorListener{}
	l.logger = logger.WithDepth(2)
	return &l
}

func (l *OmitErrorListener) OnInjectFailed(err error) {
	l.logger.Errorln(err)
}

type sliceAppender struct {
	v        reflect.Value
	elemType reflect.Type
}

func (s *sliceAppender) Set(value reflect.Value) error {
	value.Set(s.v)
	return nil
}

func (s *sliceAppender) Scan(key string, value bean.Definition) bool {
	ot := value.Type()
	// interface
	if ot.AssignableTo(s.elemType) {
		s.v = reflect.Append(s.v, value.Value())
	} else if ot.ConvertibleTo(s.elemType) {
		s.v = reflect.Append(s.v, value.Value().Convert(s.elemType))
	}
	//if s.elemType.Kind() == reflect.Interface {
	//	if ot.Implements(s.elemType) {
	//		s.v = reflect.Append(s.v, value.Value())
	//	}
	//} else if s.elemType.Kind() == reflect.Ptr {
	//	if s.elemType == value.Type() || ot.AssignableTo(s.elemType) {
	//		s.v = reflect.Append(s.v, value.Value())
	//	}
	//	//else if ot.ConvertibleTo(s.elemType) {
	//	//	s.v = reflect.Append(s.v, value.Value().Convert(s.elemType))
	//	//}
	//}

	return true
}

type mapPutter struct {
	v        reflect.Value
	elemType reflect.Type
}

func (s *mapPutter) Set(value reflect.Value) error {
	value.Set(s.v)
	return nil
}

func (s *mapPutter) Scan(key string, value bean.Definition) bool {
	ot := value.Type()
	// interface
	if ot.AssignableTo(s.elemType) {
		s.v.SetMapIndex(reflect.ValueOf(key), value.Value())
	} else if ot.ConvertibleTo(s.elemType) {
		s.v.SetMapIndex(reflect.ValueOf(key), value.Value().Convert(s.elemType))
	}
	//if s.elemType.Kind() == reflect.Interface {
	//	if ot.Implements(s.elemType) {
	//		s.v.SetMapIndex(reflect.ValueOf(key), value.Value())
	//	}
	//} else if s.elemType.Kind() == reflect.Ptr {
	//	if s.elemType == value.Type() || ot.AssignableTo(s.elemType) {
	//		s.v.SetMapIndex(reflect.ValueOf(key), value.Value())
	//	}
	//	//else if ot.ConvertibleTo(s.elemType) {
	//	//	s.v.SetMapIndex(reflect.ValueOf(key), value.Value().Convert(s.elemType))
	//	//}
	//}

	return true
}

func (injector *defaultInjector) injectInterface(c bean.Container, name string, v reflect.Value) error {
	vt := v.Type()
	if name == "" {
		name = reflection.GetTypeName(vt)
	}
	o, ok := c.GetDefinition(name)
	if ok {
		v.Set(o.Value())
		return nil
	} else {
		// 自动注入
		var matchValues []bean.Definition
		c.Scan(func(key string, value bean.Definition) bool {
			// 指定名称注册的对象直接跳过，因为在container.Get未满足，所以认定不是用户想要注入的对象
			if key != value.Name() {
				return true
			}
			ot := value.Type()
			if ot.AssignableTo(vt) {
				matchValues = append(matchValues, value)
				if len(matchValues) > 1 {
					panic("Auto Inject bean found more than 1")
				}
				return true
			}
			return true
		})
		if len(matchValues) == 1 {
			v.Set(matchValues[0].Value())
			// cache to container
			err := c.PutDefinition(reflection.GetTypeName(vt), matchValues[0])
			if err != nil {
				injector.logger.Warnln(err)
			}
			return nil
		}
	}
	return errors.New("Inject nothing, cannot find any Implementation: " + reflection.GetTypeName(vt))
}

func (injector *defaultInjector) injectSlice(c bean.Container, name string, v reflect.Value) error {
	vt := v.Type()
	if name == "" {
		name = reflection.GetSliceName(vt)
	}
	elemType := vt.Elem()
	o, ok := c.GetDefinition(name)
	if ok {
		dv := o.Value()
		n, err := util.SetOrCopySlice(v, dv, true)
		if n != dv.Len() {
			injector.logger.Infof("Set slice source have %d elements set %d elements", dv.Len(), n)
		}
		return err
	} else {
		//自动注入
		destTmp := sliceAppender{
			v:        v,
			elemType: elemType,
		}
		c.Scan(destTmp.Scan)
		destTmp.Set(v)
		if v.Len() > 0 {
			// cache to container
			bean, err := bean.CreateBeanDefinition(v.Interface())
			if err != nil {
				injector.logger.Warnln(err)
			}
			err = c.PutDefinition(reflection.GetSliceName(vt), bean)
			if err != nil {
				injector.logger.Warnln(err)
			}
			return nil
		}
	}
	return errors.New("Slice Inject nothing, cannot find any Implementation: " + reflection.GetSliceName(vt))
}

func (injector *defaultInjector) injectMap(c bean.Container, name string, v reflect.Value) error {
	vt := v.Type()
	if name == "" {
		name = reflection.GetMapName(vt)
	}
	keyType := vt.Key()
	elemType := vt.Elem()
	o, ok := c.GetDefinition(name)
	if ok {
		dv := o.Value()
		n, err := util.SetOrCopyMap(v, dv, true)
		if n != dv.Len() {
			injector.logger.Infof("Set map source have %d elements set %d elements", dv.Len(), n)
		}
		return err
	} else {
		if keyType.Kind() != reflect.String {
			return errors.New("Key type must be string. ")
		}
		//自动注入
		destTmp := mapPutter{
			v:        v,
			elemType: elemType,
		}
		c.Scan(destTmp.Scan)
		if v.Len() > 0 {
			// cache to container
			bean, err := bean.CreateBeanDefinition(v.Interface())
			if err != nil {
				injector.logger.Warnln(err)
			}
			err = c.PutDefinition(reflection.GetMapName(vt), bean)
			if err != nil {
				injector.logger.Warnln(err)
			}
			return nil
		}
	}
	return errors.New("Map Inject nothing, cannot find any Implementation: " + reflection.GetMapName(vt))
}

func (injector *defaultInjector) injectStruct(c bean.Container, name string, v reflect.Value) error {
	vt := v.Type()
	if name == "" {
		name = reflection.GetTypeName(vt)
	}
	o, ok := c.GetDefinition(name)
	if ok {
		ov := o.Value()
		if vt.Kind() == reflect.Ptr {
			v.Set(ov)
		} else {
			// 只允许注入指针类型
			err := fmt.Errorf("Inject struct: [%s] failed: value must be pointer. ", reflection.GetTypeName(vt))
			//injector.logger.Errorln(err)
			return err
			//v.Set(ov.Elem())
		}
		return nil
	}

	if injector.recursive {
		return injector.injectStructFields(c, v)
	} else {
		return errors.New("Inject nothing, cannot find any instance of  " + reflection.GetTypeName(vt))
	}
}

func (injector *defaultInjector) injectStructFields(c bean.Container, v reflect.Value) error {
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return errors.New("result must be struct ptr")
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tagAll, ok := field.Tag.Lookup(injector.tagName)
		if ok {
			tag, listeners := injector.lm.ParseListener(tagAll)
			fieldValue := v.Field(i)
			fieldType := fieldValue.Type()
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			err := injector.InjectValue(c, tag, fieldValue)
			if err != nil {
				err = fmt.Errorf("Inject failed: Field [%s: %s] error: %s\n ",
					reflection.GetTypeName(t), field.Name, err.Error())
				//injector.logger.Errorln(errStr)
				for _, l := range listeners {
					l.OnInjectFailed(err)
				}
			}
		}
	}

	return nil
}

func OptSetLogger(v xlog.Logger) Opt {
	return func(injector *defaultInjector) {
		injector.logger = v
	}
}
