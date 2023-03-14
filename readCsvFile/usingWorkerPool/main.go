package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

type City struct {
	StateCode   string
	Name        string
	CountryCode string
}

func createCity(city City) {
	time.Sleep(10 * time.Millisecond)
}

func readData(fileName string, citiesChan chan City) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	fmt.Println("CSV File opened succesfully!")

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range lines {
		citiesChan <- City{line[0], line[1], line[2]}
	}
	close(citiesChan)
}

func worker(jobs chan City) {
	for job := range jobs {
		createCity(job)
	}
}

func main() {
	startTime := time.Now()

	var cities = make(chan City, 1000)
	go readData("sample_csv.csv", cities)

	workers := 5
	job := make(chan City, 1000)
	for w := 0; w < workers; w++ {
		go worker(job)
	}

	counter := 0
	for city := range cities {
		counter++
		job <- city
	}

	fmt.Println("Records saved: ", counter)
	fmt.Println("Total time: ", time.Since(startTime))
}
