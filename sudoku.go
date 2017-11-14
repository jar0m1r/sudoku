package main

import (
	"fmt"
)

//Sudoku type represents the sudoku matrix
type sudoku [][]field

//NewSudoku generates sudoku field matrix
func newSudoku(input [][]int) sudoku {

	s := sudoku{}

	fld := field{
		value:     0,
		optionset: []int{},
		col:       []int{},
		row:       []int{},
		square:    []int{},
	}

	for i, r := range input {
		srow := []field{}
		for j, c := range r {
			f := fld
			f.value = c
			if c == 0 {
				f.optionset = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
			}
			f.pos = [2]int{i, j}
			srow = append(srow, f)
		}
		s = append(s, srow)
	}

	return s
}

//run recalculates the complete matrix, this will probably not be the most efficient. Returns true if something changed and false if no new solution
func (s sudoku) run() bool {
	for _, r := range s {
		for _, c := range r {
			c.broadcastValue(s)
		}
	}

	var resolveCnt int

	for _, r := range s {
		for index := range r {
			if r[index].resolve() {
				resolveCnt++
			}
		}
	}

	fmt.Printf("Solved %d after normal Run\n", resolveCnt)

	return resolveCnt > 0
}

func (s sudoku) guessRun() bool {
	//Better way to make full value clone
	//tempS := sudoku{}
	tempS := s
	/* 	for _, r := range s {
		srow := []field{}
		for _, c := range r {
			srow = append(srow, c)
		}
		tempS = append(tempS, srow)
	} */

	//Check dit en begrijp
	ref1 := s
	ref2 := tempS
	fmt.Printf("Original reference %p\n", &ref1)
	fmt.Printf("Original reference %p\n", &ref2)

	return true
}

func (s sudoku) getRow(f field) []*field {
	fs := []*field{}
	for i := range s[f.pos[0]] {
		fs = append(fs, &s[f.pos[0]][i])
	}
	return fs
}

func (s sudoku) getCol(f field) []*field {
	fs := []*field{}
	for _, r := range s {
		fs = append(fs, &r[f.pos[1]])
	}
	return fs
}

func (s sudoku) getSquare(f field) []*field {
	fs := []*field{}
	rowstart := f.pos[0] - f.pos[0]%3
	colstart := f.pos[1] - f.pos[1]%3

	for r := rowstart; r < rowstart+3; r++ {
		for c := colstart; c < colstart+3; c++ {
			fs = append(fs, &s[r][c])
		}
	}
	return fs
}

func (s sudoku) Print() {
	for _, r := range s {
		for _, c := range r {
			if c.value != 0 {
				fmt.Printf("%d ", c.value)
			} else {
				fmt.Printf("[%d]", len(c.optionset))
			}

		}
		fmt.Printf("\n")
	}
}
