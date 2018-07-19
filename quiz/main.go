package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	//open the csv file
	csvFile := flag.String("csv", "problems.csv", "csv contains question answer pair")
	flag.Parse()
	file, err := os.Open(*csvFile)
	if err != nil {
		exit(fmt.Sprintf("Got error in os.Open(%s) - %s", *csvFile, err.Error()))
	}

	//csv reader to read the file
	r := csv.NewReader(file) //io.Reader interface, one of the most common interface
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to read csv file"))
	}

	//parse the csv to problem struct
	problems := parseLines(lines)
	//execute quiz
	startQuiz(problems)

}
func startQuiz(problems []problem) {
	correctCounter := 0
	for i, problem := range problems {
		fmt.Printf("Problem %d: %s = \n", i+1, problem.question)
		var answer string
		fmt.Scanf("%s\n", &answer) //Scanf scans trimmed text read from standard input
		if answer == problem.answer {
			fmt.Println("Correct!")
			correctCounter++
		}
	}
	fmt.Printf("You scored %d from %d\n", correctCounter, len(problems))
}

func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))
	for i, line := range lines {
		result[i] = problem{
			strings.TrimSpace(line[0]), //strings.TrimSpace() makes it robust against invalid csv
			strings.TrimSpace(line[1]),
		}
	}
	return result
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
