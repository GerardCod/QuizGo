package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "exercises.csv", "a csv file in the formayt question,answer")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", *csvFileName))
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	if err != nil {
		exit("Failed to read the provided CSV file.")
	}

	problems := parseLines(lines)

	correct := 0

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s =\n", i+1, problem.question)
		answerChannel := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d.", correct, len(problems))
			return
		case answer := <-answerChannel:
			if answer == problem.answer {
				correct++
			}
		}
	}

}

func parseLines(lines [][]string) []problem {
	arrayOfProblems := make([]problem, len(lines))
	for i, line := range lines {
		arrayOfProblems[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return arrayOfProblems
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
