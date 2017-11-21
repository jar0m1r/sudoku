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

func (f field) broadcastValue(s sudoku) error {
	if f.value != 0 {
		fs := s.getCol(f.pos[1])
		fs = append(fs, s.getRow(f.pos[0])...)
		fs = append(fs, s.getSquare(f.pos[0], f.pos[1])...)

		for i := range fs {
			if fs[i] != &f {
				err := fs[i].blockOption("c", f.value)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (f *field) blockOption(crs string, i int) error { //crs = c(ol) r(ow) s(quare)
	if f.value == 0 {
		index := findIndex(f.optionset, i) //returns len(slice) when not found

		if index != len(f.optionset) {
			if len(f.optionset)-1 != 0 {
				(*f).optionset = append((*f).optionset[:index], (*f).optionset[index+1:]...)
			} else {
				return fmt.Errorf("No option left after blocking %d", i)
			}
		}

		switch crs {
		case "c":
			(*f).col = append((*f).col, i)
		case "r":
			(*f).row = append((*f).row, i)
		case "s":
			(*f).square = append((*f).square, i)
		}
	}
	return nil
}

func (f *field) resolve() bool {
	if f.value == 0 && len(f.optionset) == 1 {
		(*f).value = f.optionset[0]
		(*f).optionset = []int{}
		return true
	}
	return false
}

func (f *field) forceResolve(index int) {
	(*f).value = f.optionset[index]
	(*f).optionset = []int{}
}

func findIndex(data []int, v int) int {
	for index, value := range data {
		if value == v {
			return index
		}
	}
	return len(data)
}

func (f field) optionsLeft() int {
	if f.value != 0 {
		return 0
	}
	return len(f.optionset)
}

func (f field) deepClone() field {
	var fclone field
	fclone.value = f.value
	fclone.pos = f.pos
	fclone.optionset = append([]int{}, f.optionset...)
	fclone.row = append([]int{}, f.row...)
	fclone.col = append([]int{}, f.col...)
	fclone.square = append([]int{}, f.square...)
	return fclone
}
