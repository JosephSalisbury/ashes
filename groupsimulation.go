package main

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type GroupSimulation struct {
	sb, xb, yb, cellSize int
	cells, buffer        [][][]bool
}

func NewGroupSimulation(cellSize, xSize, ySize int) (*GroupSimulation, error) {
	sb := 32
	xb := xSize / cellSize
	yb := ySize / cellSize

	a := make([][][]bool, sb)
	for i := range a {
		a[i] = fillBoolCells(makeBoolCells(xb, yb))
	}

	buffer := make([][][]bool, sb)
	for i := range buffer {
		buffer[i] = makeBoolCells(xb, yb)
	}

	gs := &GroupSimulation{
		sb:       sb,
		xb:       xb,
		yb:       yb,
		cellSize: cellSize,
		cells:    a,
		buffer:   buffer,
	}

	return gs, nil
}

func makeBoolCells(xb, yb int) [][]bool {
	a := make([][]bool, xb)

	for i := range a {
		a[i] = make([]bool, yb)
	}

	return a
}

func fillBoolCells(a [][]bool) [][]bool {
	for x := 0; x < len(a); x++ {
		for y := 0; y < len(a[0]); y++ {
			a[x][y] = rand.Float32() < 0.3
		}
	}

	return a
}

func (gs *GroupSimulation) Step() error {
	for s := 0; s < len(gs.cells); s++ {
		if rand.Float32() < 0.0005 {
			gs.cells[s] = makeBoolCells(gs.xb, gs.yb)
		}
	}

	var wg sync.WaitGroup

	for s := 0; s < len(gs.cells); s++ {
		wg.Add(1)

		go func(s int) {
			defer wg.Done()

			for x := 0; x < gs.xb; x++ {
				for y := 0; y < gs.yb; y++ {
					n := gs.neighbours(s, x, y)

					gs.buffer[s][x][y] = false
					if gs.cells[s][x][y] {
						if !(n < 2 || n > 3) {
							gs.buffer[s][x][y] = true
						}
					} else {
						if n == 3 {
							gs.buffer[s][x][y] = true
						}
					}
				}
			}
		}(s)
	}

	wg.Wait()

	t := gs.cells
	gs.cells = gs.buffer
	gs.buffer = t

	return nil
}

func (gs *GroupSimulation) neighbours(s, x, y int) int {
	t := 0

	if (x != 0 && y != 0) && gs.cells[s][x-1][y-1] {
		t++
	}
	if y != 0 && gs.cells[s][x][y-1] {
		t++
	}
	if (x+1 != gs.xb && y != 0) && gs.cells[s][x+1][y-1] {
		t++
	}

	if x != 0 && gs.cells[s][x-1][y] {
		t++
	}
	if x+1 != gs.xb && gs.cells[s][x+1][y] {
		t++
	}

	if (x != 0 && y+1 != gs.yb) && gs.cells[s][x-1][y+1] {
		t++
	}
	if y+1 != gs.yb && gs.cells[s][x][y+1] {
		t++
	}
	if (x+1 != gs.xb && y+1 != gs.yb) && gs.cells[s][x+1][y+1] {
		t++
	}

	return t
}

func (gs *GroupSimulation) Render(surface *sdl.Surface) error {
	for x := 0; x < gs.xb; x++ {
		for y := 0; y < gs.yb; y++ {
			bs := []bool{}

			for s := 0; s < gs.sb; s++ {
				bs = append(bs, gs.cells[s][x][y])
			}

			c := gs.color(bs)

			rect := sdl.Rect{
				int32(x * gs.cellSize),
				int32(y * gs.cellSize),
				int32(gs.cellSize),
				int32(gs.cellSize),
			}
			surface.FillRect(&rect, c)
		}
	}

	return nil
}

func (gs *GroupSimulation) color(bs []bool) uint32 {
	var sb strings.Builder

	for _, b := range bs {
		if b {
			sb.WriteString("0")
		} else {
			sb.WriteString("1")
		}
	}

	i, _ := strconv.ParseInt(sb.String(), 2, 32)
	ix := uint32(i)

	return ix
}
