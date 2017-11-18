package main

import "fmt"

type field struct {
	value     int
	optionset []int
	col       []int
	row       []int
	square    []int
	pos       [2]int
}

func (f field) broadcastValue(s sudoku) {
	if f.value != 0 {
		fs := s.getCol(f.pos[1])
		fs = append(fs, s.getRow(f.pos[0])...)
		fs = append(fs, s.getSquare(f.pos[0], f.pos[1])...)

		for i := range fs {
			if fs[i] != &f {
				fs[i].blockOption("c", f.value)
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

func (f *field) resolve() bool {
	if f.value == 0 && len(f.optionset) == 1 {
		(*f).value = f.optionset[0]
		(*f).optionset = []int{}
		return true
	}
	return false
}

func (f *field) forceResolve(index int) bool {
	fmt.Println("got index", index)
	(*f).value = f.optionset[index]
	(*f).optionset = []int{}
	return true
}

func findIndex(data []int, v int) int {
	for index, value := range data {
		if value == v {
			return index
		}
	}
	return len(data)
}

func (f field) getOptionsLeft() int {
	return len(f.optionset)
}

func (f field) deepClone() field {
	var fclone field
	fclone.value = f.value
	fclone.pos = f.pos
	fclone.optionset = append(fclone.optionset, f.optionset...)
	fclone.row = append(fclone.row, f.row...)
	fclone.col = append(fclone.col, f.col...)
	fclone.square = append(fclone.square, f.square...)
	return fclone
}
