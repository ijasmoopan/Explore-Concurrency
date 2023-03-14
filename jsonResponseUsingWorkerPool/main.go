package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Response struct {
	Num        int    `json:"num"`
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	News       string `json:"news"`
	Link       string `json:"link"`
	Img        string `json:"img"`
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
}

type Job struct {
	number int
}

var url string = "https://xkcd.com/"
var jobs = make(chan Job, 100)
var results = make(chan Response, 100)
var resultSet []Response
var counter int = 1

func main() {

	// allocate jobs
	noOfJobs := 104
	go allocateJobs(noOfJobs)

	// get results
	done := make(chan bool)
	go getResults(done)

	// create worker pool
	noOfWorkers := 10
	createWorkerPool(noOfWorkers)

	// wait for all results to be collected
	<-done

	// convert result collection to json
	data, err := json.MarshalIndent(resultSet, "", "   ")
	if err != nil {
		log.Fatal("json error: ", err)
	}

	// write json data to file
	err = writeToFile(data)
	if err != nil {
		log.Fatal("writeToFile error: ", err)
	}
	log.Println(len(resultSet), " responses saved to 'response.json' file.")
}

func writeToFile(data []byte) error {
	fmt.Println("WRITE TO FILE")
	file, err := os.Create("response.json")
	if err != nil && !os.IsExist(err) {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func getResults(done chan bool) {
	fmt.Println("GET RESULTS")
	for result := range results {
		resultSet = append(resultSet, result)
	}
	done <- true
}

func worker(wg *sync.WaitGroup) {
	fmt.Println("WORKER")
	for job := range jobs {
		fmt.Println(counter)
		counter++
		result, err := fetch(job.number)
		if err != nil {
			log.Printf("Error in fetch: %d - %v\n", job.number, err)
		}
		results <- result
	}
	wg.Done()
}

func createWorkerPool(noOfWorkers int) {
	fmt.Println("CREATE WORKER POOL")
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
	close(results)
}

func allocateJobs(noOfJobs int) {
	fmt.Println("ALLOCATE JOBS")
	for i := 1; i <= noOfJobs; i++ {
		jobs <- Job{i}
	}
	close(jobs)
}

func fetch(jobNo int) (Response, error) {

	fmt.Println("FETCH")

	currentUrl := fmt.Sprintf("%s%d/info.0.json", url, jobNo)

	var result Response

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest("GET", currentUrl, nil)
	if err != nil {
		return result, err
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return result, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}
	fmt.Println("result id: ", result.Num)

	return result, nil
}
