package util

import (
	"github.com/xfali/xlog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HandlerSignal(logger xlog.Logger, closers ...func() error) (err error) {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(100 * time.Millisecond)
			xlog.Infof("Got a signal %s, closing...", si.String())
			go func() {
				for i := range closers {
					cErr := closers[i]()
					if cErr != nil {
						logger.Errorln(cErr)
						err = cErr
					}
				}
			}()
			time.Sleep(3 * time.Second)
			xlog.Infof("------ Process exited ------")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
