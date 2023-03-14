package main

import (
	"fmt"
	"sync"
	"time"
)

type JobInput struct {
	StartTime time.Time
	EndTime   time.Time
}

type JobOutput struct{}

type Response struct {
	Message string
}

func startWorker(jobsInputChan <-chan JobInput, jobsOutputChan chan<- JobOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	for jobInput := range jobsInputChan {
		jobsOutputChan <- buisnessFunctionalityJob(jobInput)
	}
}

func buisnessFunctionality(jobsInputChan chan<- JobInput) {
	jobInputs := []JobInput{JobInput{}, JobInput{}}
	for _, jobInput := range jobInputs {
		jobsInputChan <- jobInput
	}
	close(jobsInputChan)
}

func buisnessFunctionalityJob(jobInput JobInput) JobOutput {
	fmt.Println("Job is executing...")
	return JobOutput{}
}

func responseGenerator(responses []JobOutput) Response {
	return Response{"Jobs finished!"}
}

func main() {

	noOfWorkers := 3
	jobsChan := make(chan JobInput, 10)
	resultsChan := make(chan JobOutput, 10)

	wg := sync.WaitGroup{}

	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go startWorker(jobsChan, resultsChan, &wg)
	}

	go buisnessFunctionality(jobsChan)

	var wgResp sync.WaitGroup
	var responses []JobOutput
	wgResp.Add(1)
	go func() {
		defer wgResp.Done()
		for result := range resultsChan {
			responses = append(responses, result)
		}
	}()

	wg.Wait()
	close(resultsChan)
	wgResp.Wait()

	fmt.Println(responseGenerator(responses).Message)
}
