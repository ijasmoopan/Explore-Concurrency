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

func main() {
	startTime := time.Now()

	csvFile, err := os.Open("sample_csv.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	counter := 0
	for _, line := range lines {
		counter++
		createCity(City{line[0], line[1], line[2]})
	}

	fmt.Println("Records saved: ", counter)
	fmt.Println("Total time: ", time.Since(startTime))
}
