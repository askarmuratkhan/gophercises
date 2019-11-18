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
	csvFileName := flag.String("csv", "problems.csv", "csv file with format \"problem, answer\"")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	file, err := os.Open(*csvFileName)
	if err != nil {
		Exit(fmt.Sprintf("Failed to open file: %s\n", *csvFileName))
	}
	read := csv.NewReader(file)
	lines, err := read.ReadAll()
	problems := ParseProblems(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			//fmt.Printf("\nCorrect answer %d out of %d: ", correct, len(problems))
			fmt.Println()
			//return
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}

	}

	fmt.Printf("Correct answer %d out of %d: ", correct, len(problems))

}

func ParseProblems(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

type problem struct {
	q string
	a string
}
