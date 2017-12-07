package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
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
	start := time.Now()
	solution := solve(s)
	elapsed := time.Since(start)

	fmt.Printf("Solution is \n%v\nin %s", solution.Print(), elapsed)

}

func solve(s sudoku) sudoku {
	done := make(chan bool)
	queue := make(chan sudoku)
	solved := make(chan sudoku)

	guessqueue := []sudoku{}

	doRun(done, queue, solved, s)

	for {
		select {
		case ns := <-queue:
			guessqueue = append(guessqueue, ns)
		case <-done:
			if len(guessqueue) > 0 {
				lastS := guessqueue[len(guessqueue)-1]
				guessqueue = guessqueue[:len(guessqueue)-1]
				gs := lastS.guess()
				doRun(done, queue, solved, gs...)
				continue
			}
			return sudoku{}
		case sol := <-solved:
			return sol
		}
	}
}

func doRun(done chan bool, queue chan sudoku, solved chan sudoku, suds ...sudoku) {

	var wg sync.WaitGroup
	wg.Add(len(suds))

	for _, su := range suds {
		go func(s sudoku) {
			if err := s.run(); err == nil {
				if s.isValid() {
					if s.isSolved() {
						solved <- s
						wg.Done()
						return
					}
					queue <- s
				}
			}
			wg.Done()
		}(su)
	}

	go func() {
		wg.Wait()
		done <- true
	}()

}
