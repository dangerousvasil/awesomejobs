package aweorchestrator

import (
	"awesomejobs/awejob"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"sync"
)

type AweOrchestrator struct {
	ctx     context.Context
	mutex   *sync.RWMutex
	JobMap  map[uuid.UUID]*awejob.AweJob
	handler func(ctx context.Context, parameters json.RawMessage) error
	wg      *sync.WaitGroup
}

func NewAweOrchestrator(ctx context.Context, handler func(ctx context.Context, parameters json.RawMessage) error) (*AweOrchestrator, error) {
	ork := AweOrchestrator{
		ctx:     ctx,
		handler: handler,
		mutex:   new(sync.RWMutex),
		wg:      new(sync.WaitGroup),
		JobMap:  map[uuid.UUID]*awejob.AweJob{},
	}
	go ork.watchdog()
	return &ork, nil
}

func (o *AweOrchestrator) SearchJobs(name string) ([]*awejob.AweJob, error) {
	o.mutex.RLock()
	jobs := []*awejob.AweJob{}
	for u := range o.JobMap {
		if o.JobMap[u].GetName() == name {
			jobs = append(jobs, o.JobMap[u])
		}
	}
	o.mutex.RUnlock()
	return jobs, nil
}

func (o *AweOrchestrator) GetJob(jobid uuid.UUID) (*awejob.AweJob, error) {
	o.mutex.RLock()
	job, ok := o.JobMap[jobid]
	o.mutex.RUnlock()
	if ok {
		return job, nil
	}
	return nil, nil
}

func (o *AweOrchestrator) Stop() {
	for u := range o.JobMap {
		o.JobMap[u].Stop()
	}
}
