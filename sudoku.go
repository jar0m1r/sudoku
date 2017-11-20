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

func (s sudoku) solve(c chan sudoku) {
	for {
		r := s.run()
		if r == -1 {
			fmt.Println("This sudoku cannot be solved")
			c <- nil
		} else if r == 0 {
			fmt.Println("Another run finished but didn't resolve anything")
			break
		} else {
			fmt.Println("Another run finished")
		}
	}

	if !s.isSolved() {
		s = s.guessRun()
	} else {
		c <- s
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

	fmt.Printf("Solved %d \n", resolveCnt)

	return resolveCnt
}

func (s sudoku) guessRun() sudoku {
	minoptions := 2
	c := make(chan sudoku)
Outerloop:
	for x := minoptions; x < 9; x++ {
		for i := 0; i < 9; i++ {
			for _, f := range s.getCol(i) {
				if len(f.optionset) == x {
					fmt.Println("least options", x)
					go guessBranch(0, f.pos, s, c)
					go guessBranch(1, f.pos, s, c)
					break Outerloop
				}
			}
		}
		minoptions++
	}

	return <-c
}

func guessBranch(index int, pos [2]int, s sudoku, c chan sudoku) {
	sclone := s.deepClone()
	row := pos[0]
	col := pos[1]
	sclone[row][col].forceResolve(index)

	d := make(chan sudoku)
	go sclone.solve(d)
	c <- (<-d)
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
