package appcontext

import (
	"errors"
	"github.com/xfali/xlog"
	"reflect"
	"sync"
)

const (
	defaultEventBufferSize = 4096
)

var eventType = reflect.TypeOf((*ApplicationEvent)(nil)).Elem()

type defaultEventProcessor struct {
	logger xlog.Logger

	listeners    []ApplicationEventListener
	listenerLock sync.Mutex

	eventBufSize int
	eventChan    chan ApplicationEvent

	consumerListenerFac func() ApplicationEventConsumerListener

	stopChan   chan struct{}
	finishChan chan struct{}
	closeOnce  sync.Once
}

type EventProcessorOpt func(processor *defaultEventProcessor)

func NewEventProcessor(opts ...EventProcessorOpt) *defaultEventProcessor {
	ret := &defaultEventProcessor{
		logger:       xlog.GetLogger(),
		eventBufSize: defaultEventBufferSize,
		//consumerListenerFac: defaultConsumerListenerFac,
	}

	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

type dummyEventProc struct{}

func NewDisableEventProcessor() *dummyEventProc {
	return &dummyEventProc{}
}

func (p *dummyEventProc) NotifyEvent(e ApplicationEvent) error {
	panic("Application event process: Disabled")
}

func (p *dummyEventProc) PublishEvent(e ApplicationEvent) error {
	panic("Application event process: Disabled")
}

func (p *dummyEventProc) AddListeners(listeners ...interface{}) {
	panic("Application event process: Disabled")
}

func (p *dummyEventProc) Start() error {
	return nil
}

func (p *dummyEventProc) Close() error {
	return nil
}

type PayloadEventListener struct {
	invokers []ConsumerInvoker
}

type PayloadApplicationEvent struct {
	BaseApplicationEvent
	payload interface{}
}

func NewPayloadApplicationEvent(payload interface{}) *PayloadApplicationEvent {
	if payload == nil {
		return nil
	}
	e := PayloadApplicationEvent{}
	e.ResetOccurredTime()
	e.payload = payload
	return &e
}

func (l *PayloadEventListener) OnApplicationEvent(e ApplicationEvent) {
	if len(l.invokers) > 0 {
		if pe, ok := e.(*PayloadApplicationEvent); ok {
			for _, invoker := range l.invokers {
				invoker.Invoke(pe.payload)
			}
		}
	}
}

type ConsumerInvoker interface {
	// 消费
	Invoke(data interface{}) bool

	// 检查consumer是否符合类型要求
	ResolveConsumer(consumer interface{}) error
}

type consumerInvoker struct {
	et reflect.Type
	fv reflect.Value
}

type eventInvoker struct {
	consumerInvoker
}

func (invoker *eventInvoker) ResolveConsumer(consumer interface{}) error {
	t := reflect.TypeOf(consumer)
	if t.Kind() != reflect.Func {
		return errors.New("Param is not a function. ")
	}

	if t.NumIn() != 1 {
		return errors.New("Param is not match, expect func(ApplicationEvent). ")
	}

	et := t.In(0)
	if !et.AssignableTo(eventType) {
		return errors.New("Param is not match, function param must Implements ApplicationEvent. ")
	}

	invoker.et = et
	invoker.fv = reflect.ValueOf(consumer)
	return nil
}

type payloadInvoker struct {
	consumerInvoker
}

func (invoker *payloadInvoker) ResolveConsumer(consumer interface{}) error {
	t := reflect.TypeOf(consumer)
	if t.Kind() != reflect.Func {
		return errors.New("Param is not a function. ")
	}

	if t.NumIn() != 1 {
		return errors.New("Param is not match, expect func(ApplicationEvent). ")
	}

	et := t.In(0)
	invoker.et = et
	invoker.fv = reflect.ValueOf(consumer)
	return nil
}
