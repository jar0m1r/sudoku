package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var input = [][]int{}

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

	var mutex sync.Mutex
	var counter int

	for i := 0; i < 4; i++ {
		go func(num int, in chan sudoku, out chan sudoku) {
			for s := range in {
				mutex.Lock()
				counter++
				mutex.Unlock()
				if err := s.runCycle(); err == nil {
					if s.isValid() {
						out <- s
					} else {
						//fmt.Printf("Invalid solution on Routine %d \n", num)
					}
				} else {
					//fmt.Printf("Error received on Routine %d, error : %s\n", num, err)
				}
				mutex.Lock()
				counter--
				mutex.Unlock()
			}
		}(i, cRun, cGuess)
	}

	for i := 0; i < 10000; i++ {
		go func(num int, in chan sudoku, out chan sudoku, solution chan sudoku) {
			for s := range in {
				mutex.Lock()
				counter++
				mutex.Unlock()
				if !s.isSolved() {
					for _, clone := range s.guess() {
						out <- clone
					}
				} else {
					solution <- s
				}
				mutex.Lock()
				counter--
				mutex.Unlock()
			}
		}(i, cGuess, cRun, cSolution)
	}

	go func() {
		time.Sleep(time.Millisecond * 100)
		for {
			if counter == 0 {
				close(cRun)
				close(cGuess)
				close(cSolution)
				break
			}
		}
	}()

	cRun <- s

	nSolution := 0
	for solution := range cSolution {
		nSolution++
		fmt.Printf("\nSolution %d\n %s \n", nSolution, solution.Print())
		if nSolution > 0 {
			break
		}
	}

}
