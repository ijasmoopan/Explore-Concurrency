package main

import (
	"fmt"
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

	for _, jobInput := range jobInputs {
		func() {
			jobOutputs = append(jobOutputs, buisnessFunctionalityJob(jobInput))
			fmt.Println("Executing job...")
		}()
	}
	return responseGenerator(jobOutputs)
}

func buisnessFunctionalityJob(jobInput JobInput) JobOutput {
	return JobOutput{}
}

func responseGenerator(repsonses []JobOutput) Response {
	return Response{Message: "Jobs finished!"}
}

func main() {
	fmt.Println(buisnessFunctionality().Message)
}
