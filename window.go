package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"fmt"
)

const (
	FRAME_CLOCKS = 70224 // number of t-cycle per frame
)

type Window struct {
	width int
	height int
	screen *sdl.Window
	renderer *sdl.Renderer
	texture *sdl.Texture
	fps int
	running bool
}

func (window *Window) init() {
	window.width = 160
	window.height = 144
	window.running = true
	sdl.Init(sdl.INIT_EVERYTHING)
	var err error
	window.screen, err = sdl.CreateWindow("DreamGB", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 256, 256, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	
	window.renderer, err = sdl.CreateRenderer(window.screen, -1, sdl.RENDERER_ACCELERATED)
    if err != nil {
    	panic(err)
	}	
}

func (window *Window) loop() {
	for window.running {
		// tick components
		if config.stepmode {
			var input string
			fmt.Print("> ")
			fmt.Scanf("%s", &input)
			if input == "" || input == "next" || input == "n" {
				cpu.step()
			}

		} else {
			for cpu.cycles < FRAME_CLOCKS {
				cpu.cycles++
			}
		}
    	window.renderer.Clear()
    	window.renderer.SetDrawColor(255, 255, 55, 255)
		window.renderer.Present()
    	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				window.running = false
			case *sdl.KeyboardEvent:
				if t.Keysym.Sym == sdl.K_ESCAPE {
					window.running = false
				}
			}
		}
    	sdl.Delay(1000 / 60)
	}
	defer window.screen.Destroy()
	defer window.renderer.Destroy()
}