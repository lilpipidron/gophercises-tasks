package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	time2 "time"
)

func csvOpen(filename string) (quests [][]string, err error) {
	file, err := os.Open(filename)
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

func timer(dur time2.Duration, breakQuiz chan int) {
	time := time2.NewTimer(dur * time2.Second)
	<-time.C
	fmt.Println("Time Out")
	breakQuiz <- 1
}

func main() {
	var (
		count  int
		answer string
	)

	breakQuiz := make(chan int, 1)
	quests, err := csvOpen("problems.csv")
	if err != nil {
		log.Fatal(err)
	}

	go timer(5, breakQuiz)
	for _, i := range quests {
		if len(breakQuiz) != 0 {
			break
		}
		fmt.Println(i[0])

		_, err := fmt.Scan(&answer)
		if err != nil {
			log.Fatal(err)
		}

		if answer == i[1] {
			count++
		}
	}
	fmt.Println(count, len(quests))
}
