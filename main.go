package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	csvFilename string
	timeLimit   int
)

type problem struct {
	q string
	a string
}

func readArgs() {
	flag.StringVar(&csvFilename, "csv", "problems.csv", "a csv of questions")
	flag.IntVar(&timeLimit, "limit", 30, "Time limit for quiz is 30seconds")

}

func openFile() []problem {
	file, err := os.Open(csvFilename)
	if err != nil {
		exit(fmt.Sprintln("couldnt open the file"))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("could not read file")
	}

	quiz := make([]problem, len(lines))
	for i, line := range lines {
		quiz[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return quiz
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	readArgs()
	flag.Parse()
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	problems := openFile()

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You got %d out of %d correct.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You got %d out of %d correct.\n", correct, len(problems))
}

