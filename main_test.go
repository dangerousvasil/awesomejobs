package main

import (
	"awesomejobs/aweorchestrator"
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	jobHandler := func(ctx context.Context, parameters json.RawMessage) error {
		log.Println("parameters ", string(parameters))
		for i := 0; i < 7; i++ {
			time.Sleep(time.Second)
			log.Println("run ", i)
			log.Println(ctx.Err())
		}
		log.Println("job super end")
		return nil
	}
	aweOrk, err := aweorchestrator.NewAweOrchestrator(ctx, jobHandler)
	if err != nil {
		panic(err)
	}
	name := "name of job"
	var parameters json.RawMessage = []byte("my json parameters")

	jobid, err := aweOrk.RunJob(ctx, name, parameters)
	if err != nil {
		panic(err)
	}

	job1s, err := aweOrk.SearchJobs(name)
	if err != nil {
		panic(err)
	}
	for i := range job1s {
		log.Println(job1s[i].GetUUID(), job1s[i].GetName())
	}

	job2, err := aweOrk.GetJob(jobid.GetUUID())
	if err != nil {
		return
	}
	log.Println(job2.GetUUID(), job2.GetName())

	time.Sleep(time.Second)
	time.Sleep(time.Second)

	log.Println("---catch end---")
	cancel()

	aweOrk.Wait()
	log.Println("happy end")
}
