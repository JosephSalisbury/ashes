package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type GroupSimulation struct {
	sb, xb, yb, cellSize int
	cells                [][][]bool
}

func NewGroupSimulation(cellSize, xSize, ySize int) (*GroupSimulation, error) {
	sb := 32
	xb := xSize / cellSize
	yb := ySize / cellSize

	a := make([][][]bool, sb)

	for i := range a {
		a[i] = fillCells(makeCells(xb, yb))
	}

	gs := &GroupSimulation{
		sb:       sb,
		xb:       xb,
		yb:       yb,
		cellSize: cellSize,
		cells:    a,
	}

	return gs, nil
}

func makeCells(xb, yb int) [][]bool {
	a := make([][]bool, xb)

	for i := range a {
		a[i] = make([]bool, yb)
	}

	return a
}

func fillCells(a [][]bool) [][]bool {
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
			gs.cells[s] = makeCells(gs.xb, gs.yb)
		}
	}

	for s := 0; s < len(gs.cells); s++ {
		// TODO: Remove array copy.
		t := makeCells(gs.xb, gs.yb)

		for x := 0; x < gs.xb; x++ {
			for y := 0; y < gs.yb; y++ {
				n := gs.neighbours(s, x, y)

				if gs.cells[s][x][y] {
					if n < 2 || n > 3 {
						t[x][y] = false
					} else {
						t[x][y] = true
					}
				} else {
					if n == 3 {
						t[x][y] = true
					}
				}
			}
		}

		gs.cells[s] = t
	}

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
