package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type (
	Problem struct {
		Question string
		Answer   string
	}
)

func main() {
	var file string
	var limit int
	var shuffle bool
	flag.StringVar(&file, "file", "problems.csv", "the path to file with questions")
	flag.IntVar(&limit, "limit", 30, "the timer for answering questions in seconds")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle questions")
	flag.Parse()

	// parse time limit

	f, err := os.OpenFile(file, os.O_RDONLY, 0400)
	defer f.Close()
	if err != nil {
		log.Fatalf("operation OpenFile error: %v", err)
	}

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("operation ReadAll for csv error: %v", err)
	}

	problems := parseProblems(rows)

	if shuffle {
		problems = shuffleProblems(problems)
	}

	timer := time.NewTimer(time.Duration(limit) * time.Second)
	var score int
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.Question)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scan(&answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTimes runs out. Your score is %d/%d\n", score, len(rows))
			return
		case answer := <-answerCh:
			if strings.ToLower(answer) == p.Answer {
				score++
			}
		}
	}

	fmt.Printf("Your score is %d/%d\n", score, len(rows))
}

func parseProblems(lines [][]string) []Problem {
	problems := make([]Problem, len(lines))
	for i, l := range lines {
		problems[i] = Problem{
			Question: l[0],
			Answer:   strings.ToLower(strings.TrimSpace(l[1])),
		}
	}
	return problems
}

func shuffleProblems(problems []Problem) []Problem {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
	return problems
}
