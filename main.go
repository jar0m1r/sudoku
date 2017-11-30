package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	cDone := make(chan bool)
	cGuess := make(chan sudoku)

	for i := 0; i < 4; i++ {
		go func(num int, in chan sudoku, done chan bool, out chan sudoku) {
			for s := range in {
				if err := s.runCycle(); err == nil {
					fmt.Printf("\nRun cycle result\n %s \n----\n", s.Print())
					out <- s
				} else {
					fmt.Printf("Error received on Routine %d, error : %s\n", num, err)
				}
			}
		}(i, cRun, cDone, cGuess)
	}

	for i := 0; i < 10000; i++ {
		go func(num int, in chan sudoku, done chan bool, out chan sudoku) {
			for s := range in {
				fmt.Printf("\nReceived on Guess channel %d\n %v \n --- \n", num, s.Print())
				if !s.isSolved() {
					for _, clone := range s.guess() {
						out <- clone
					}
				} else {
					fmt.Printf("\nSOLVED on channel %d!!\n %s \n----\n", num, s.Print())
					done <- true
					break
				}
			}
			/* 		time.Sleep(time.Second * 2)
			   		fmt.Println("Closing guess channel")
			   		close(in) */
		}(i, cGuess, cDone, cRun)
	}

	cRun <- s

mainloop:
	for {
		select {
		case (<-cDone):
			//time.Sleep(time.Second * 3)
			fmt.Println("All done..")
			break mainloop
		default:
			//fmt.Println("Waiting for done signal")
		}
	}

}
