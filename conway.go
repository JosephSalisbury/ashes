package main

var (
	black uint32 = 0x00000000
	white uint32 = 0xffffffff
)

type Simulation struct {
	size  int
	array [][]bool
}

func NewSimulation(size, x, y int) (*Simulation, error) {
	xp := x / size
	yp := y / size

	a := make([][]bool, xp)

	for i := range a {
		a[i] = make([]bool, yp)
	}

	s := &Simulation{
		size:  size,
		array: a,
	}

	for i := 0; i < xp; i++ {
		for j := 0; j < yp; j++ {
			s.array[i][j] = randBool()
		}
	}

	return s, nil
}

func (s *Simulation) Step() error {
	t := make([][]bool, len(s.array))
	for x := range t {
		t[x] = make([]bool, len(s.array[0]))
	}

	for i := 0; i < len(s.array); i++ {
		for j := 0; j < len(s.array[0]); j++ {
			n := getNumLiveNeighbours(s.array, i, j)

			if s.array[i][j] {
				if n < 2 || n > 3 {
					t[i][j] = false
				} else {
					t[i][j] = true
				}
			} else {
				if n == 3 {
					t[i][j] = true
				}
			}
		}
	}

	s.array = t

	return nil
}

func (s *Simulation) GetCell(x, y int) bool {
	return s.array[x][y]
}

func getNumLiveNeighbours(a [][]bool, x, y int) int {
	t := 0

	// top left
	if (x != 0 && y != 0) && a[x-1][y-1] {
		t++
	}

	// top
	if y != 0 && a[x][y-1] {
		t++
	}

	// top right
	if (x+1 != len(a) && y != 0) && a[x+1][y-1] {
		t++
	}

	// left
	if x != 0 && a[x-1][y] {
		t++
	}

	// right
	if x+1 != len(a) && a[x+1][y] {
		t++
	}

	// bottom left
	if (x != 0 && y+1 != len(a[0])) && a[x-1][y+1] {
		t++
	}

	// bottom
	if y+1 != len(a[0]) && a[x][y+1] {
		t++
	}

	// bottom right
	if (x+1 != len(a) && y+1 != len(a[0])) && a[x+1][y+1] {
		t++
	}

	return t
}
