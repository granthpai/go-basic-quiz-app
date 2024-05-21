package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func problemPuller(fileName string) ([]problem, error) {
	// open the file
	fObj, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error in opening %s file: %s", fileName, err.Error())
	}
	defer fObj.Close()

	// create new reader
	csvR := csv.NewReader(fObj)
	// read the file
	cLines, err := csvR.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error in reading data in csv format from %s file: %s", fileName, err.Error())
	}

	// call parseProblem func
	return parseProblem(cLines), nil
}

func main() {
	// input name of file
	fName := flag.String("f", "quiz.csv", "path of csv file")
	// set timer
	timer := flag.Int("t", 30, "timer of the quiz")
	flag.Parse()

	// pull problems
	problems, err := problemPuller(*fName)
	if err != nil {
		exit(fmt.Sprintf("something went wrong: %s", err.Error()))
	}

	// create variable to count correct answers
	correctAns := 0
	// using duration of timer we initialize timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

	// loop through problems, print questions, and accept answers
problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s = ", i+1, p.q)

		go func() {
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()

		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ansC:
			if strings.TrimSpace(iAns) == strings.TrimSpace(p.a) {
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}

	// calculate and print result
	fmt.Printf("Your result is %d out of %d\n", correctAns, len(problems))
	fmt.Printf("Press enter to exit")
	fmt.Scanln()
}

func parseProblem(lines [][]string) []problem {
	// go over lines and parse them
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
