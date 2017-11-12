package main

import (
	"fmt"
)

type field struct {
	value     int
	optionset []int
	col       []int
	row       []int
	square    []int
	pos       [2]int
}

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

func (s sudoku) initRun() {
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

	fmt.Printf("Solved %d after init Run\n", resolveCnt)
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

func guessRun() bool {
	return true
}

func (f field) broadcastValue(s sudoku) {
	if f.value != 0 {
		//col broadcast
		for _, r := range s {
			r[f.pos[1]].blockOption("c", f.value)
		}
		//row broadcast
		for index := range s[f.pos[0]] {
			s[f.pos[0]][index].blockOption("r", f.value)
		}
		//square broadcast
		rowstart := f.pos[0] - f.pos[0]%3
		colstart := f.pos[1] - f.pos[1]%3
		for r := rowstart; r < rowstart+3; r++ {
			for c := colstart; c < colstart+3; c++ {
				s[r][c].blockOption("s", f.value)
			}
		}
	}
}

func (f *field) blockOption(crs string, i int) { //crs = c(ol) r(ow) s(quare)
	if f.value == 0 {
		switch crs {
		case "c":
			(*f).col = append((*f).col, i)
		case "r":
			(*f).row = append((*f).row, i)
		case "s":
			(*f).square = append((*f).square, i)
		}

		index := findIndex(f.optionset, i)

		if index != len(f.optionset) {
			(*f).optionset = append((*f).optionset[:index], (*f).optionset[index+1:]...)
		}
	}
}

func findIndex(data []int, v int) int {
	for index, value := range data {
		if value == v {
			return index
		}
	}
	return len(data)
}

func (f *field) resolve() bool {
	if f.value == 0 && len(f.optionset) == 1 {
		(*f).value = f.optionset[0]
		(*f).optionset = []int{}
		return true
	}
	return false
}

func (s sudoku) Print() {
	for _, r := range s {
		for _, c := range r {
			fmt.Printf("%d ", c.value)
		}
		fmt.Printf("\n")
	}
}
