package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type OptimisedGroupSimulation struct {
	xb, yb, cellSize int
	cells            [][]uint32
}

func NewOptimisedGroupSimulation(cellSize, xSize, ySize int) (*OptimisedGroupSimulation, error) {
	xb := xSize / cellSize
	yb := ySize / cellSize

	a := fillUint32Cells(makeUint32Cells(xb, yb))

	ogs := &OptimisedGroupSimulation{
		xb:       xb,
		yb:       yb,
		cellSize: cellSize,
		cells:    a,
	}

	return ogs, nil
}

func makeUint32Cells(xb, yb int) [][]uint32 {
	a := make([][]uint32, xb)

	for i := range a {
		a[i] = make([]uint32, yb)
	}

	return a
}

func fillUint32Cells(a [][]uint32) [][]uint32 {
	for x := 0; x < len(a); x++ {
		for y := 0; y < len(a[0]); y++ {
			a[x][y] = rand.Uint32()
		}
	}

	return a
}

func (ogs *OptimisedGroupSimulation) Step() error {
	a := makeUint32Cells(ogs.xb, ogs.yb)

	for s := 0; s < 32; s++ {
		for x := 0; x < ogs.xb; x++ {
			for y := 0; y < ogs.yb; y++ {
				live := get(s, ogs.cells[x][y])
				n := ogs.neighbours(s, x, y)

				if live {
					if n < 2 || n > 3 {
						a[x][y] = set(s, ogs.cells[x][y], false)
					} else {
						a[x][y] = set(s, ogs.cells[x][y], true)
					}
				} else {
					if n == 3 {
						a[x][y] = set(s, ogs.cells[x][y], true)
					}
				}
			}
		}
	}

	ogs.cells = a

	return nil
}

func (ogs *OptimisedGroupSimulation) neighbours(s, x, y int) int {
	t := 0

	if (x != 0 && y != 0) && get(s, ogs.cells[x-1][y-1]) {
		t++
	}
	if y != 0 && get(s, ogs.cells[x][y-1]) {
		t++
	}
	if (x+1 != ogs.xb && y != 0) && get(s, ogs.cells[x+1][y-1]) {
		t++
	}

	if x != 0 && get(s, ogs.cells[x-1][y]) {
		t++
	}
	if x+1 != ogs.xb && get(s, ogs.cells[x+1][y]) {
		t++
	}

	if (x != 0 && y+1 != ogs.yb) && get(s, ogs.cells[x-1][y+1]) {
		t++
	}
	if y+1 != ogs.yb && get(s, ogs.cells[x][y+1]) {
		t++
	}
	if (x+1 != ogs.xb && y+1 != ogs.yb) && get(s, ogs.cells[x+1][y+1]) {
		t++
	}

	return t
}

func get(s int, a uint32) bool {
	if strings.Split(fmt.Sprintf("%032b", a), "")[s] == "1" {
		return true
	}

	return false
}

func set(s int, a uint32, alive bool) uint32 {
	t := strings.Split(fmt.Sprintf("%032b", a), "")

	if alive {
		t[s] = "1"
	} else {
		t[s] = "0"
	}

	i, _ := strconv.ParseInt(strings.Join(t, ""), 2, 32)
	return uint32(i)
}

func (ogs *OptimisedGroupSimulation) Render(surface *sdl.Surface) error {
	for x := 0; x < ogs.xb; x++ {
		for y := 0; y < ogs.yb; y++ {
			rect := sdl.Rect{
				int32(x * ogs.cellSize),
				int32(y * ogs.cellSize),
				int32(ogs.cellSize),
				int32(ogs.cellSize),
			}
			surface.FillRect(&rect, ogs.cells[x][y])
		}
	}

	return nil
}
