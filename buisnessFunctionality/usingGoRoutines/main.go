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

func buisnessFunctionality() Response {
	jobInputs := []JobInput{JobInput{}, JobInput{}, JobInput{}}
	jobOutputs := []JobOutput{}

	wg := sync.WaitGroup{}

	for _, jobInput := range jobInputs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			jobOutputs = append(jobOutputs, buisnessFunctionalityJob(jobInput))
			fmt.Println("Executing Job...")
		}()
	}
	wg.Wait()
	return responseGenerator(jobOutputs)
}

func buisnessFunctionalityJob(jobInput JobInput) JobOutput {
	return JobOutput{}
}

func responseGenerator(responses []JobOutput) Response {
	return Response{"Jobs finished"}
}

func main() {
	fmt.Println(buisnessFunctionality().Message)
}
