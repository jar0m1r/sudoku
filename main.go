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
	//s.initRun()

	for s.run() {
		fmt.Println("Another run started")
	}

	s.Print()

	s.guessRun()
}
