package aweorchestrator

import (
	"time"
)

func (o *AweOrchestrator) Wait() {
	chanfirst := make(chan bool, 1)
	go func() {
		o.wg.Wait()
		chanfirst <- true
	}()
	timer := time.NewTimer(10 * time.Second)
	select {
	case <-timer.C:
	case <-chanfirst:
	}

}

func (o *AweOrchestrator) watchdog() {
	select {
	case <-o.ctx.Done():
		o.Stop()
	}
}
