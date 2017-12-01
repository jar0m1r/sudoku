package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var input = [][]int{}

func init() {
	nf, err := os.Create("log.txt")
	if err != nil {
		fmt.Println("err")
	}
	log.SetOutput(nf)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		b := []rune(scanner.Text())
		irow := []int{}

		for i := 0; i < len(b); i++ {
			if b[i] != 32 {
				if v, err := strconv.Atoi(string(b[i])); err == nil {
					irow = append(irow, v)

				} else {
					irow = append(irow, 0)
				}
			}
		}
		input = append(input, irow)

		if len(input) >= 9 {
			break
		}
	}

	s := newSudoku(input)

	solve(s)

}

func solve(s sudoku) {

	cRun := make(chan sudoku)
	cGuess := make(chan sudoku)
	cSolution := make(chan sudoku)

	for i := 0; i < 100; i++ {
		go func(num int, in chan sudoku, out chan sudoku) {
			for s := range in {
				log.Printf("run routine %d received from channel in\n", num)
				if err := s.runCycle(); err == nil {
					if s.isValid() {
						log.Printf("run routine %d sending to channel out\n", num)
						out <- s
					}
				}
				log.Printf("run routine %d done with loop\n", num)
			}
		}(i, cRun, cGuess)
	}

	for i := 0; i < 100; i++ {
		go func(num int, in chan sudoku, out chan sudoku, solution chan sudoku) {
			for s := range in {
				log.Printf("guess routine %d received from channel in\n", num)
				if !s.isSolved() {
					for _, clone := range s.guess() {
						log.Printf("quess routine %d sending to channel out\n", num)
						out <- clone
					}
				} else {
					log.Println("Solution found..")
					log.Printf("quess routine %d sending to channel solution\n", num)
					solution <- s
				}
				log.Printf("guess routine %d done with loop\n", num)
			}
		}(i, cGuess, cRun, cSolution)
	}

	cRun <- s

	nSolution := 0
	for solution := range cSolution {
		nSolution++
		log.Printf("\nSolution %d\n %s \n", nSolution, solution.Print())
	}

}
