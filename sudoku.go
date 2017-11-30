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

func (s sudoku) runCycle() error {
	var nRun int
	for {
		r := s.run()
		if r == -1 {
			//there is at least one field without possible solutions
			return fmt.Errorf("This sudoku is faulty")
		} else if r == 0 {
			//Run finished but didn't resolve anything
			return nil
		} else {
			//Another run finished and solved some new
		}
		nRun++
	}
}

//run recalculates the complete matrix, this will probably not be the most efficient. Returns true if something changed and false if no new solution
func (s sudoku) run() int {
	for _, r := range s {
		for _, c := range r {
			err := c.broadcastValue(s)
			if err != nil {
				return -1
			}
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

	return resolveCnt
}

func (s sudoku) guess() []sudoku {
	alloptions := []map[int][]*field{
		s.getColOptions(),
		s.getRowOptions(),
		s.getSquareOptions(),
	}

	combinedOptions := map[int][]*field{}
	fieldsDone := map[*field]bool{}

	for x := 2; x < 10; x++ {
		for _, optionsMap := range alloptions {
			if fs, ok := optionsMap[x]; ok {
				for _, f := range fs {
					if _, ok := fieldsDone[f]; !ok {
						combinedOptions[x] = append(combinedOptions[x], f)
						fieldsDone[f] = true
					}
				}
			}
		}
	}

	result := []sudoku{}
	for i := 2; i < 10; i++ {
		if v, ok := combinedOptions[i]; ok {
			for j := 0; j < i; j++ {
				result = append(result, guessBranch(j, v[0].pos, s))
			}
			break
		}
	}
	return result
}

func guessBranch(index int, pos [2]int, s sudoku) sudoku {
	sclone := s.deepClone()
	row := pos[0]
	col := pos[1]
	sclone[row][col].forceResolve(index)

	return sclone
}

func (s sudoku) getRow(row int) []*field {
	fs := []*field{}
	for i := range s[row] {
		fs = append(fs, &s[row][i])
	}
	return fs
}

func (s sudoku) getCol(col int) []*field {
	fs := []*field{}
	for _, r := range s {
		fs = append(fs, &r[col])
	}
	return fs
}

func (s sudoku) getSquare(row, col int) []*field {
	fs := []*field{}
	rowstart := row - row%3
	colstart := col - col%3

	for r := rowstart; r < rowstart+3; r++ {
		for c := colstart; c < colstart+3; c++ {
			fs = append(fs, &s[r][c])
		}
	}
	return fs
}

func (s sudoku) getColOptions() map[int][]*field {
	var colOptionsMap = map[int][]*field{}
	leastSoFar := 100
	for i := 0; i < 9; i++ {
		optionsleft := 0
		var tempMap = map[int][]*field{}
		for _, f := range s.getCol(i) {
			if n := f.optionsLeft(); n > 0 {
				optionsleft = optionsleft + n
				tempMap[n] = append(tempMap[n], &s[f.pos[0]][f.pos[1]])
			}
		}
		if optionsleft > 0 && optionsleft < leastSoFar {
			leastSoFar = optionsleft
			colOptionsMap = tempMap
		}
	}
	return colOptionsMap
}

func (s sudoku) getRowOptions() map[int][]*field {
	var rowOptionsMap = map[int][]*field{}
	leastSoFar := 100
	for i := 0; i < 9; i++ {
		optionsleft := 0
		var tempMap = map[int][]*field{}
		for _, f := range s.getRow(i) {
			if n := f.optionsLeft(); n > 0 {
				optionsleft = optionsleft + n
				tempMap[n] = append(tempMap[n], &s[f.pos[0]][f.pos[1]])
			}
		}
		if optionsleft > 0 && optionsleft < leastSoFar {
			leastSoFar = optionsleft
			rowOptionsMap = tempMap
		}
	}
	return rowOptionsMap
}

func (s sudoku) getSquareOptions() map[int][]*field {
	var squareOptionsMap = map[int][]*field{}
	leastSoFar := 100
	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			optionsleft := 0
			var tempMap = map[int][]*field{}
			for _, f := range s.getSquare(i, j) {
				if n := f.optionsLeft(); n > 0 {
					optionsleft = optionsleft + n
					tempMap[n] = append(tempMap[n], &s[f.pos[0]][f.pos[1]])
				}
			}
			if optionsleft > 0 && optionsleft < leastSoFar {
				leastSoFar = optionsleft
				squareOptionsMap = tempMap
			}
		}
	}
	return squareOptionsMap
}

func (s sudoku) Print() string {
	var result string
	for _, r := range s {
		for _, c := range r {
			if c.value != 0 {
				result += fmt.Sprintf("%d ", c.value)
			} else {
				result += fmt.Sprintf("%s", "_ ") //[%d] len(c.optionset)
			}

		}
		result += fmt.Sprintf("\n")
	}
	return result
}

func (s sudoku) isSolved() bool {
	for i := 0; i < 9; i++ {
		for _, f := range s.getRow(i) {
			if f.value == 0 {
				return false
			}
		}
	}
	return true
}

func (s sudoku) isValid() bool {
	for r := 0; r < 9; r++ {
		if !uniqueSet(s.getRow(r)) {
			return false
		}
		for c := 0; c < 9; c++ {
			if !uniqueSet(s.getCol(c)) {
				return false
			}
		}
	}

	for r := 0; r < 9; r += 3 {
		for c := 0; c < 9; c += 3 {
			if !uniqueSet(s.getSquare(r, c)) {
				return false
			}
		}
	}
	return true
}

func uniqueSet(fs []*field) bool {
	valuesMap := map[int]bool{}
	for _, f := range fs {
		v := f.value
		if v != 0 {
			if _, ok := valuesMap[v]; ok {
				return false
			}
			valuesMap[v] = true
		}
	}
	return true
}

func (s sudoku) deepClone() sudoku {
	var sclone sudoku = [][]field{}
	for index, r := range s {
		sclone = append(sclone, []field{})
		for i := range r {
			sclone[index] = append(sclone[index], r[i].deepClone())
		}
	}
	return sclone
}
