package aweorchestrator

import (
	"awesomejobs/awejob"
	"context"
	"encoding/json"
	"errors"
)

func (o *AweOrchestrator) RunJob(ctx context.Context, name string, parameters json.RawMessage) (*awejob.AweJob, error) {
	job, err := awejob.NewAweJob(ctx, name, parameters, o.handler)
	if err != nil {
		return nil, err
	}
	o.mutex.Lock()
	if _, ok := o.JobMap[job.GetUUID()]; !ok {
		o.JobMap[job.GetUUID()] = job
	} else {
		err = errors.New("job with same uuid exist")
	}
	o.mutex.Unlock()
	o.wg.Add(1)
	go o.run(job)

	return job, err
}

func (o *AweOrchestrator) run(job *awejob.AweJob) {

	var err error
	var errchan = make(chan error, 1)

	errchan <- job.Run()
	select {
	case <-o.ctx.Done():
		err = o.ctx.Err()
	case err = <-errchan:
	}

	job.SetError(err)
	o.wg.Done()
}
