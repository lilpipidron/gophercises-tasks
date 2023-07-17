package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func openCSV(filename string) (quests [][]string, err error) {
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

func runQuiz(breakQuizChan chan bool, answersChan chan string, quests [][]string) error {
	var answer string

	for _, q := range quests {

		fmt.Println(q[0])
		_, err := fmt.Scan(&answer)
		if err != nil {
			return err
		}

		select {
		case <-breakQuizChan:
			close(answersChan)
			return nil
		default:
			answersChan <- answer
		}

	}
	close(answersChan)
	return nil
}

func checkAns(answersChan chan string, quests [][]string) int {
	var (
		count, point int
	)
	for data := range answersChan {
		if data == quests[point][1] {
			count++
		}
		point++
	}
	return count
}

func main() {
	breakQuizChan := make(chan bool, 1)
	quests, err := openCSV("problems.csv")
	if err != nil {
		log.Fatalf("Unable to open CSV file: %v", err)
	}
	answersChan := make(chan string, len(quests))

	go func() {
		timer := time.NewTimer(5 * time.Second)
		<-timer.C
		fmt.Println("Time Out")
		breakQuizChan <- true
		close(breakQuizChan)
	}()

	err = runQuiz(breakQuizChan, answersChan, quests)

	if err != nil {
		log.Fatalf("Unable to run quiz: %v", err)
	}

	count := checkAns(answersChan, quests)

	fmt.Println(count, len(quests))
}
