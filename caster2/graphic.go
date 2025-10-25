package main

import (
	"encoding/binary"
	"fmt"
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
	TEX_WIDTH     = 64
	TEX_HEIGHT    = 64
)

var worldMap = [MAP_WIDTH][MAP_HEIGHT]int{
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 7, 7, 7, 7, 7, 7, 7, 7},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0, 0, 7},
	{4, 0, 4, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 7, 7, 0, 7, 7, 7, 7, 7},
	{4, 0, 5, 0, 0, 0, 0, 5, 0, 5, 0, 5, 0, 5, 0, 5, 7, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 6, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 0, 0, 0, 8},
	{4, 0, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 8, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 0, 0, 0, 8},
	{4, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 5, 7, 0, 0, 0, 7, 7, 7, 1},
	{4, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 0, 5, 5, 5, 5, 7, 7, 7, 7, 7, 7, 7, 1},
	{6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
	{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
	{6, 6, 6, 6, 6, 6, 0, 6, 6, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6},
	{4, 4, 4, 4, 4, 4, 0, 4, 4, 4, 6, 0, 6, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 2, 0, 0, 5, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2},
	{4, 0, 6, 0, 6, 0, 0, 0, 0, 4, 6, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 2},
	{4, 0, 0, 5, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2},
	{4, 0, 6, 0, 6, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 5, 0, 0, 2, 0, 0, 0, 2},
	{4, 0, 0, 0, 0, 0, 0, 0, 0, 4, 6, 0, 6, 2, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2},
	{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 1, 1, 1, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3},
}

type Graphic struct {
	window          *sdl.Window
	surface         *sdl.Surface
	pos, dir, plane vec2f
	time, oldTime   uint64 // time of current frame and time of previous frame
	frameTime       float64
	buffer          [SCREEN_HEIGHT][SCREEN_WIDTH]uint32
	texture         [8][TEX_WIDTH * TEX_HEIGHT]uint32
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

	g.surface, err = g.window.GetSurface()
	if err != nil {
		return fmt.Errorf("failed to get surface: %s", err)
	}
	g.surface.FillRect(nil, sdl.MapRGB(g.surface.Format, 0x10, 0x10, 0x10))
	g.window.UpdateSurface()

	g.pos = vec2f{22, 11.5}  // x and y start position
	g.dir = vec2f{-1, 0}     // initial direction vector
	g.plane = vec2f{0, 0.66} // the 2d raycaster version of camera plane

	for x := range uint32(TEX_WIDTH) {
		for y := range uint32(TEX_HEIGHT) {
			xorcolor := (x * 256 / TEX_WIDTH) ^ (y * 256 / TEX_HEIGHT)
			//int xcolor = x * 256 / textureWidth;
			ycolor := y * 256 / TEX_HEIGHT
			xycolor := y*128/TEX_HEIGHT + x*128/TEX_WIDTH

			if x != y && x != TEX_WIDTH-y {
				g.texture[0][TEX_WIDTH*y+x] = 65536 * 254
			} else {
				g.texture[0][TEX_WIDTH*y+x] = 0
			}

			g.texture[1][TEX_WIDTH*y+x] = xycolor + 256*xycolor + 65536*xycolor
			g.texture[2][TEX_WIDTH*y+x] = 256*xycolor + 65536*xycolor
			g.texture[3][TEX_WIDTH*y+x] = xorcolor + 256*xorcolor + 65536*xorcolor
			g.texture[4][TEX_WIDTH*y+x] = 256 * xorcolor
			if x%16 != 0 && y%16 != 0 {
				g.texture[5][TEX_WIDTH*y+x] = 65536 * 192
			} else {
				g.texture[5][TEX_WIDTH*y+x] = 0
			}
			g.texture[6][TEX_WIDTH*y+x] = 65536 * ycolor
			g.texture[7][TEX_WIDTH*y+x] = 128 + 256*128 + 65536*128
		}
	}

	return nil
}

func (g *Graphic) CleanUp() {
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

		pitch := 100

		// calculate lowest and highest pixel to fill in current stripe
		drawStart := max(-lineHeight/2+SCREEN_HEIGHT/2+pitch, 0)
		drawEnd := lineHeight/2 + SCREEN_HEIGHT/2 + pitch
		if drawEnd >= SCREEN_HEIGHT {
			drawEnd = SCREEN_HEIGHT - 1
		}

		// texturing calculations
		textureNum := worldMap[ipos.x][ipos.y] - 1 // 1 subtracted from it so that texture 0 can be used!

		// calculate value of wallX
		wallX := 0.0 // where exactly the wall was hit
		if !side {
			wallX = g.pos.y + perpWallDist*rayDir.y
		} else {
			wallX = g.pos.x + perpWallDist*rayDir.x
		}
		wallX -= math.Floor(wallX)

		// x coordinate on the texture
		texX := int(wallX * TEX_WIDTH)

		if !side && rayDir.x > 0 {
			texX = TEX_WIDTH - texX - 1
		}
		if side && rayDir.y < 0 {
			texX = TEX_WIDTH - texX - 1
		}

		// How much to increase the texture coordinate per screen pixel
		texStep := 1.0 * TEX_HEIGHT / float64(lineHeight)
		// Starting texture coordinate
		texPos := (float64(drawStart-pitch) - SCREEN_HEIGHT/2 + float64(lineHeight)/2) * texStep
		for y := drawStart; y < drawEnd; y++ {
			// Cast the texture coordinate to integer, and mask with (texHeight - 1) in case of overflow
			textureY := int(texPos) & (TEX_HEIGHT - 1)
			texPos += texStep
			c := g.texture[textureNum][TEX_HEIGHT*textureY+texX]
			// make color darker for y-sides: R, G and B byte each divided through two with a "shift" and an "and"
			if side {
				c = (c >> 1) & 8355711
			}
			g.buffer[y][x] = c
		}
	}
}

func (g *Graphic) Render() {
	offsetX := (WINDOW_WIDTH - SCREEN_WIDTH) / 2
	offsetY := (WINDOW_HEIGHT - SCREEN_HEIGHT) / 2

	g.surface.Lock()
	defer g.surface.Unlock()
	pixels := g.surface.Pixels()
	for y := range int(SCREEN_HEIGHT) {
		for x := range int(SCREEN_WIDTH) {
			index := (y+offsetY)*int(g.surface.Pitch) + (x+offsetX)*4
			binary.LittleEndian.PutUint32(pixels[index:index+4], g.buffer[y][x])
		}
	}

	for y := range SCREEN_HEIGHT {
		for x := range SCREEN_WIDTH {
			g.buffer[y][x] = 0
		}
	}
	g.window.UpdateSurface()
}
