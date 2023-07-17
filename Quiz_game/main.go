package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func readCSV(filename string) (quests [][]string, err error) {
	file, err := os.Open(filename)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	if err != nil {
		return
	}
	reader := csv.NewReader(file)
	quests, err = reader.ReadAll()
	if err != nil {
		return
	}
	return
}

func runQuiz(score *int, quests [][]string) error {
	var (
		answer string
	)

	for _, q := range quests {
		fmt.Println(q[0])
		_, err := fmt.Scan(&answer)
		if err != nil {
			return err
		}
		if q[1] == answer {
			*score++
		}
	}
	return nil
}

func main() {
	quests, err := readCSV("problems.csv")
	if err != nil {
		log.Fatalf("Unable to open CSV file: %v", err)
	}
	timeoutCh := time.After(5 * time.Second)

	var score int
	go func() {
		err = runQuiz(&score, quests)
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatalf("Unable to run quiz: %v", err)
	}

	select {
	case <-timeoutCh:
		fmt.Println("Timeout")
	}

	fmt.Println(score, len(quests))
}
