package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindowName = "ashes"

	WindowSizeX = 1200
	WindowSizeY = 800

	CellSize = 5

	FramesPerSecond = 60
	StepsPerSecond  = 30
)

func init() {
	rand.Seed(time.Now().UnixNano())
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

	var simulation Simulation
	{
		s, err := NewGroupSimulation(CellSize, WindowSizeX, WindowSizeY)
		if err != nil {
			log.Fatalf("could not create simulation: %v", err)
		}

		simulation = s
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

		surface.FillRect(nil, 0)

		if err := simulation.Step(); err != nil {
			log.Fatalf("could not step simulation: %v", err)
		}

		if err := simulation.Render(surface); err != nil {
			log.Fatalf("could not render simulation: %v", err)
		}

		window.UpdateSurface()

		sdl.Delay(30)
	}
}
