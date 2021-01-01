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

	CellSize       = 20
	StepsPerSecond = 20
	Debug          = true
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
		gs, err := NewGroupSimulation(CellSize, WindowSizeX, WindowSizeY)
		if err != nil {
			log.Fatalf("could not create group simulation: %v", err)
		}

		simulation = gs
	}
	// {
	// 	ogs, err := NewOptimisedGroupSimulation(CellSize, WindowSizeX, WindowSizeY)
	// 	if err != nil {
	// 		log.Fatalf("could not create optimised group simulation: %v", err)
	// 	}
	//
	// 	simulation = ogs
	// }

	maxDuration := time.Second / StepsPerSecond
	lastStepTime := time.Now()

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

		now := time.Now()
		var stepStart time.Time
		var renderStart time.Time
		var renderEnd time.Time

		if now.Sub(lastStepTime) > maxDuration {
			surface.FillRect(nil, 0)

			stepStart = time.Now()

			if err := simulation.Step(); err != nil {
				log.Fatalf("could not step simulation: %v", err)
			}

			renderStart = time.Now()

			if err := simulation.Render(surface); err != nil {
				log.Fatalf("could not render simulation: %v", err)
			}

			renderEnd = time.Now()

			lastStepTime = now

			window.UpdateSurface()
		}

		duration := time.Now().Sub(now)

		var waitTime uint32
		if duration < maxDuration {
			waitTime = uint32((maxDuration - duration).Milliseconds())
		}

		if Debug {
			log.Printf("frame: %v (step: %v, render: %v) (max: %v), wait: %v ms", duration, renderStart.Sub(stepStart), renderEnd.Sub(renderStart), maxDuration, waitTime)
		}

		if waitTime > 0 {
			sdl.Delay(waitTime)
		}
	}
}
