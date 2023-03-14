package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Job struct {
	id int
}
type Result struct {
	job  Job
	file string
}

var alphabet = 'a'

var jobs = make(chan Job)
var results = make(chan Result)

func generateFileName(alphabet rune) string {
	fileName := fmt.Sprintf("1-%v", string(alphabet))
	time.Sleep(1 * time.Second)
	fmt.Printf("%s-> ", string(alphabet))

	return fileName
}

func worker(wg *sync.WaitGroup, mutex *sync.Mutex) {
	for job := range jobs {
		mutex.Lock()
		result := Result{job, generateFileName(alphabet)}
		results <- result
		alphabet++
		mutex.Unlock()
	}
	wg.Done()
}

func createWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg, &mutex)
	}
	wg.Wait()
	close(results)
}

func allocate(noOfJobs int) {
	for i := 0; i < noOfJobs; i++ {
		jobs <- Job{i}
	}
	close(jobs)
}

func result(done chan bool, m *sync.Mutex) {
	for result := range results {
		m.Lock()
		file, err := os.Create(fmt.Sprintf("GeneratedFiles/%s", result.file))
		if err != nil {
			log.Fatal(err)
		}
		m.Unlock()
		defer file.Close()
	}
	done <- true
}

func main() {
	noOfJobs := 26
	go allocate(noOfJobs)

	done := make(chan bool)
	var m sync.Mutex

	go result(done, &m)

	noOfWorkers := 4
	createWorkerPool(noOfWorkers)

	<-done
	fmt.Println(noOfJobs, "file generated")
}
