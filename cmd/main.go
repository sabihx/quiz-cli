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
	var csvFilename *string = flag.String("csv", "problems.csv", "a csv file of questions in the format 'question,answer'")
	var timeLimit *int = flag.Int("timelimit", 30, "time limit for quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the given CSV file.")
	}

	var problems []problem = parseLines(lines)

	var startTime time.Time = time.Now()
	var endTime time.Time
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correctCount := 0

	problemsLoop:
		for i, p := range problems {
			fmt.Printf("Problem #%d: %s = ", i+1, p.q)
			answerCh := make(chan string)
			go func() {
				var answer string
				fmt.Scanf("%s\n", &answer)
				answerCh <- answer
			}()

			select {
			case <- timer.C:
				fmt.Println()
				break problemsLoop
			case answer := <- answerCh:
				if answer == p.a {
					correctCount++
				}
			}
		}
	
	endTime = time.Now()
	fmt.Printf("You scored %d out of %d in %.2f seconds", correctCount, len(problems), float64(endTime.Sub(startTime)) / float64(time.Second))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string 
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}