package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func main() {
	file, err := os.Open("problems.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	timer := time.NewTimer(30 * time.Second)

	totalQuestions := len(records)
	correctAnswers := 0
loop:
	for _, record := range records {
		fmt.Println("Question:", record[0])
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanln(&answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Println("Time's up!")
			break loop
		case answer := <-answerCh:
			if answer == record[1] {
				correctAnswers++
			}
		}

	}
	fmt.Printf("You answered %d/%d questions correctly.\n", correctAnswers, totalQuestions)
}
