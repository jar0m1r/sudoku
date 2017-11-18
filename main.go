package main

import (
	"bufio"
	"os"
	"strconv"
)

var input = [][]int{}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	c := make(chan sudoku)

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

	s.solve(c)

	(<-c).Print()

}
