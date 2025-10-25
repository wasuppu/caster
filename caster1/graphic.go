package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WINDOW_WIDTH  = 1280
	WINDOW_HEIGHT = 720
	SCREEN_WIDTH  = 640
	SCREEN_HEIGHT = 480
	MAP_WIDTH     = 24
	MAP_HEIGHT    = 24
)

var worldMap = [MAP_WIDTH][MAP_HEIGHT]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

type Graphic struct {
	window          *sdl.Window
	renderer        *sdl.Renderer
	pos, dir, plane vec2f
	time, oldTime   uint64 // time of current frame and time of previous frame
	frameTime       float64
}

func (g *Graphic) Init() (err error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return fmt.Errorf("failed to initializing SDL: %s", err)
	}

	g.window, err = sdl.CreateWindow("raycast", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_BORDERLESS|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		return fmt.Errorf("failed to create SDL window: %s", err)
	}

	g.renderer, err = sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %s", err)
	}

	g.pos = vec2f{22, 12}    // x and y start position
	g.dir = vec2f{-1, 0}     // initial direction vector
	g.plane = vec2f{0, 0.66} // the 2d raycaster version of camera plane

	return nil
}

func (g *Graphic) CleanUp() {
	g.renderer.Destroy()
	g.window.Destroy()
	sdl.Quit()
}

func (g *Graphic) Done() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			fmt.Println("\nQuitting...")
			return true
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				return true
			}
		}
	}
	return false
}

func (g *Graphic) HandleKeys() {
	// speed modifiers
	moveSpeed := g.frameTime * 5.0 // the constant value is in squares/second
	rotSpeed := g.frameTime * 3.0  // the constant value is in radians/second

	sdl.PumpEvents()
	keystate := sdl.GetKeyboardState()

	// move forward if no wall in front of you
	if keystate[sdl.SCANCODE_UP] != 0 {
		if worldMap[int(g.pos.x+g.dir.x*moveSpeed)][int(g.pos.y)] == 0 {
			g.pos.x += g.dir.x * moveSpeed
		}

		if worldMap[int(g.pos.x)][int(g.pos.y+g.dir.y*moveSpeed)] == 0 {
			g.pos.y += g.dir.y * moveSpeed
		}
	}

	// move backwards if no wall behind you
	if keystate[sdl.SCANCODE_DOWN] != 0 {
		if worldMap[int(g.pos.x-g.dir.x*moveSpeed)][int(g.pos.y)] == 0 {
			g.pos.x -= g.dir.x * moveSpeed
		}

		if worldMap[int(g.pos.x)][int(g.pos.y-g.dir.y*moveSpeed)] == 0 {
			g.pos.y -= g.dir.y * moveSpeed
		}
	}

	// rotate to the right
	if keystate[sdl.SCANCODE_RIGHT] != 0 {
		// both camera direction and camera plane must be rotated
		oldDir := g.dir
		g.dir.x = g.dir.x*math.Cos(-rotSpeed) - g.dir.y*math.Sin(-rotSpeed)
		g.dir.y = oldDir.x*math.Sin(-rotSpeed) + g.dir.y*math.Cos(-rotSpeed)
		oldPlane := g.plane
		g.plane.x = g.plane.x*math.Cos(-rotSpeed) - g.plane.y*math.Sin(-rotSpeed)
		g.plane.y = oldPlane.x*math.Sin(-rotSpeed) + g.plane.y*math.Cos(-rotSpeed)
	}

	// rotate to the left
	if keystate[sdl.SCANCODE_LEFT] != 0 {
		// both camera direction and camera plane must be rotated
		oldDir := g.dir
		g.dir.x = g.dir.x*math.Cos(rotSpeed) - g.dir.y*math.Sin(rotSpeed)
		g.dir.y = oldDir.x*math.Sin(rotSpeed) + g.dir.y*math.Cos(rotSpeed)
		oldPlane := g.plane
		g.plane.x = g.plane.x*math.Cos(rotSpeed) - g.plane.y*math.Sin(rotSpeed)
		g.plane.y = oldPlane.x*math.Sin(rotSpeed) + g.plane.y*math.Cos(rotSpeed)
	}
}

func (g *Graphic) SetFrameTime() {
	// timing for input and FPS counter
	g.oldTime = g.time
	g.time = sdl.GetTicks64()
	g.frameTime = float64(g.time-g.oldTime) / 1000.0 // frametime is the time this frame has taken, in seconds
}

func (g *Graphic) Calculate() {
	// WALL CASTING
	for x := range SCREEN_WIDTH {
		// calculate ray position and direction
		cameraX := 2*float64(x)/SCREEN_WIDTH - 1 // x-coordinate in camera space

		rayDir := g.dir.add(g.plane.muln(cameraX))

		// which box of the map we're in
		ipos := vec2i{int(g.pos.x), int(g.pos.y)}

		// length of ray from one x or y-side to next x or y-side
		deltaDist := vec2f{}
		if rayDir.x == 0 {
			deltaDist.x = 1e30
		} else {
			deltaDist.x = math.Abs(1 / rayDir.x)
		}

		if rayDir.y == 0 {
			deltaDist.y = 1e30
		} else {
			deltaDist.y = math.Abs(1 / rayDir.y)
		}

		// length of ray from current position to next x or y-side
		sideDist := vec2f{}
		// what direction to step in x or y-direction (either +1 or -1)
		step := vec2f{}

		// calculate step and initial sideDist
		if rayDir.x < 0 {
			step.x = -1
			sideDist.x = (g.pos.x - float64(ipos.x)) * deltaDist.x
		} else {
			step.x = 1
			sideDist.x = (float64(ipos.x) + 1.0 - g.pos.x) * deltaDist.x
		}

		if rayDir.y < 0 {
			step.y = -1
			sideDist.y = (g.pos.y - float64(ipos.y)) * deltaDist.y
		} else {
			step.y = 1
			sideDist.y = (float64(ipos.y) + 1.0 - g.pos.y) * deltaDist.y
		}

		hit := false  // was there a wall hit?
		side := false // was a NS or a EW wall hit?

		// perform DDA
		for !hit {
			// jump to next map square, either in x-direction, or in y-direction
			if sideDist.x < sideDist.y {
				sideDist.x += deltaDist.x
				ipos.x = int(float64(ipos.x) + step.x)
				side = false
			} else {
				sideDist.y += deltaDist.y
				ipos.y = int(float64(ipos.y) + step.y)
				side = true
			}

			// Check if ray has hit a wall
			if worldMap[ipos.x][ipos.y] > 0 {
				hit = true
			}
		}

		perpWallDist := 0.0
		// Calculate distance of perpendicular ray (Euclidean distance would give fisheye effect!)
		if !side {
			perpWallDist = sideDist.x - deltaDist.x
		} else {
			perpWallDist = sideDist.y - deltaDist.y
		}

		// Calculate height of line to draw on screen
		lineHeight := int(SCREEN_HEIGHT / perpWallDist)

		// calculate lowest and highest pixel to fill in current stripe
		drawStart := max(-lineHeight/2+SCREEN_HEIGHT/2, 0)
		drawEnd := lineHeight/2 + SCREEN_HEIGHT/2
		if drawEnd >= SCREEN_HEIGHT {
			drawEnd = SCREEN_HEIGHT - 1
		}

		c := color.RGBA{}

		switch worldMap[ipos.x][ipos.y] {
		case 1:
			c = red
		case 2:
			c = green
		case 3:
			c = blue
		case 4:
			c = white
		default:
			c = yellow
		}
		if side {
			c = color.RGBA{c.R / 2, c.G / 2, c.B / 2, c.A}
		}

		offsetX := (WINDOW_WIDTH - SCREEN_WIDTH) / 2
		offsetY := (WINDOW_HEIGHT - SCREEN_HEIGHT) / 2

		g.verline(x+offsetX, drawStart+offsetY, drawEnd+offsetY, c)
	}
}

func (g *Graphic) verline(x, y1, y2 int, c color.RGBA) {
	g.renderer.SetDrawColor(c.R, c.G, c.B, c.A)
	g.renderer.DrawLine(int32(x), int32(y1), int32(x), int32(y2))
}

func (g *Graphic) Render() {
	g.renderer.Present()
	g.renderer.SetDrawColor(black.R, black.G, black.B, black.A)
	g.renderer.Clear()
}
