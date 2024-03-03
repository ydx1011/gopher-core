package appcontext

import (
	"github.com/xfali/xlog"
	"sync"
)

const (
	defaultEventBufferSize = 4096
)

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
