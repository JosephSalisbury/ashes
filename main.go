package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindowName = "ashes"

	WindowSizeX = 1200
	WindowSizeY = 800

	CellSize = 7
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randUint8() uint8 {
	return uint8(rand.Intn(255))
}

func randBool() bool {
	return rand.Float32() < 0.25
}

func getColor(bs []bool) uint32 {
	s := ""

	for _, b := range bs {
		x := "1"
		if b {
			x = "0"
		}

		s = s + x
	}

	i, _ := strconv.ParseInt(s, 2, 32)
	ix := uint32(i)

	return ix
}

func Render(simulations []*Simulation, surface *sdl.Surface) error {
	x := WindowSizeX / CellSize
	y := WindowSizeY / CellSize

	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			cs := []bool{}
			for _, s := range simulations {
				cs = append(cs, s.GetCell(i, j))
			}

			c := getColor(cs)

			rect := sdl.Rect{int32(i * CellSize), int32(j * CellSize), int32(CellSize), int32(CellSize)}
			surface.FillRect(&rect, c)
		}
	}

	return nil
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("could not initialise SDL: %v", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		WindowName,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		WindowSizeX,
		WindowSizeY,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Fatalf("could not create window: %v", err)
	}
	defer window.Destroy()

	simulations := []*Simulation{}

	for i := 0; i < 32; i++ {
		s, err := NewSimulation(CellSize, WindowSizeX, WindowSizeY)
		if err != nil {
			log.Fatalf("could not create simulation %v: %v", i, err)
		}

		simulations = append(simulations, s)
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		surface, err := window.GetSurface()
		if err != nil {
			log.Fatalf("could not get surface: %v", err)
		}

		for i := 0; i < len(simulations); i++ {
			if rand.Float32() < 0.0001 {
				s, err := NewSimulation(CellSize, WindowSizeX, WindowSizeY)
				if err != nil {
					log.Fatalf("could not recreate simulation %v: %v", i, err)
				}

				simulations[i] = s
			}
		}

		var wg sync.WaitGroup
		for _, s := range simulations {
			wg.Add(1)
			go func(s *Simulation) {
				defer wg.Done()
				if err := s.Step(); err != nil {
					log.Fatalf("could not step simulation: %v", err)
				}
			}(s)
		}
		wg.Wait()

		surface.FillRect(nil, 0)

		if err := Render(simulations, surface); err != nil {
			log.Fatalf("could not render simulation: %v", err)
		}

		window.UpdateSurface()

		sdl.Delay(30)
	}
}
