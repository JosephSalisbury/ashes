package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Simulation interface {
	Step() error
	Render(*sdl.Surface) error
}
