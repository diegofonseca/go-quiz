package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	goshuffle "github.com/secnot/goshuffle"
)

func main() {

	csvFileName := flag.String("csv", "problem.csv", "A csv file containing 'question,answer'")
	timeLimit := flag.Int("limit", 30, "Number of seconds to complete the test")
	shuffle := flag.Bool("shuffle", false, "Number of seconds to complete the test")
	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err!= nil {
        exit(fmt.Sprintf("Fail to open csv file: %v\n", csvFileName))
    }

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err!= nil {
		exit(fmt.Sprintf("Fail to parse csv file: %v\n", csvFileName))
	}

	problems := parseLines(lines)

	if *shuffle {
		goshuffle.Shuffle(problems, goshuffle.NewRandSource())
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

	problemloop:
		for index, problem := range problems {
			fmt.Printf("Problem #%d: %s = \n", index+1, problem.question)
			answerCh := make(chan string)

			go func() {
				var answer string
				fmt.Scanf("%s", &answer)
				answerCh <- strings.TrimSpace(answer)
			}()

			select {
				case <-timer.C:
					fmt.Println()
					break problemloop
				case answer := <- answerCh:
					if strings.TrimSpace(problem.answer) == answer {
						correct++
					}
			}
		}

	fmt.Printf("You scored %d out of %d.\n", correct, len(lines))

}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))
    for i, line := range lines {
        problems[i] = problem{
            question: line[0],
            answer:   line[1],
        }
    }
	return problems
}

type problem struct {
	question string
	answer  string
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}
