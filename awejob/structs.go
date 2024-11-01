package awejob

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"log"
	"time"
)

type AweJob struct {
	jobId   uuid.UUID
	name    string
	handle  func(context.Context, json.RawMessage) error
	params  json.RawMessage
	running bool
	tmStart time.Time
	logs    []string
	ctx     context.Context
	cancel  context.CancelFunc
	error   error
}

func (j *AweJob) GetName() string {
	return j.name
}

func (j *AweJob) GetUUID() uuid.UUID {
	return j.jobId
}

func (j *AweJob) Run() error {
	j.running = true
	var err error
	var errchan = make(chan error, 1)
	errchan <- j.handle(j.ctx, j.params)
	select {
	case <-j.ctx.Done():
		err = j.ctx.Err()
	case err = <-errchan:
	}
	log.Println("exit jrun")
	j.running = false
	return err
}

func (j *AweJob) GetLogs() []string {
	return j.logs
}
func (j *AweJob) AddLog(log string) {
	j.logs = append(j.logs, log)
}

func (j *AweJob) Context() context.Context {
	return j.ctx
}
func (j *AweJob) Stop() {
	j.cancel()
}

func (j *AweJob) SetRunning(b bool) {
	j.running = b
}

func (j *AweJob) SetError(err error) {
	j.error = err
}

func NewAweJob(ctx context.Context, name string, parameters json.RawMessage, handler func(ctx context.Context, parameters json.RawMessage) error) (*AweJob, error) {
	v7, err := uuid.NewV7()
	if err != nil {
		return nil, nil
	}

	awej := AweJob{
		jobId:   v7,
		name:    name,
		handle:  handler,
		params:  parameters,
		running: false,
	}

	awej.ctx, awej.cancel = context.WithCancel(ctx)

	return &awej, err
}
